package relay

import (
	"bytes"
	"errors"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"socksv/network/smux"
	"socksv/protocol"
	"socksv/protocol/socks5"
)

const (
	Version1    byte = 0x01
	CmdConnect  byte = 0x01
	EncryptNone byte = 0x00
	EncryptAES  byte = 0x01

	StatusSuccess            = 0x00
	StatusUnSupportedVersion = 0x01
	StatusUnsupportedCmd     = 0x02
	StatusUnSupportedEncrypt = 0x03
	StatusTargetUnreachable  = 0x04
	ProtocolID               = 1
	PingProtocolID           = 2
)
const readLen = 1024 * 2

type Request struct {
	Version byte
	Cmd     byte
	Encrypt byte
	Length  byte
	Data    []byte
}

func (req *Request) Read(reader io.Reader) error {
	he := make([]byte, 4)
	if _, err := io.ReadFull(reader, he); err != nil {
		return err
	}
	req.Version = he[0]
	req.Cmd = he[1]
	req.Encrypt = he[2]
	req.Length = he[3]
	c := make([]byte, req.Length)
	if _, err := io.ReadFull(reader, c); err != nil {
		return err
	}
	req.Data = c
	return nil
}

func (h Request) Bytes() []byte {
	return append([]byte{h.Version, h.Cmd, h.Encrypt, h.Length}, h.Data...)
}

type Reply struct {
	//status:0-success
	Version byte
	Status  byte
	Length  byte
	Data    []byte
}

func (rep *Reply) Read(reader io.Reader) error {
	he := make([]byte, 3)
	if _, err := io.ReadFull(reader, he); err != nil {
		return err
	}
	rep.Version = he[0]
	rep.Status = he[1]
	rep.Length = he[2]
	c := make([]byte, rep.Length)
	if _, err := io.ReadFull(reader, c); err != nil {
		return err
	}
	rep.Data = c
	return nil
}

func (h Reply) Bytes() []byte {
	return append([]byte{h.Version, h.Status, h.Length}, h.Data...)
}

func writeRequest(writer io.Writer, cmd byte, data []byte) error {
	header := Request{
		Version: Version1,
		Cmd:     cmd,
		Encrypt: EncryptAES,
		Length:  byte(len(data)),
		Data:    data,
	}
	if _, err := writer.Write(header.Bytes()); err != nil {
		return err
	}
	return nil
}
func writeReply(writer io.Writer, status byte, data []byte) error {
	header := Reply{
		Version: Version1,
		Status:  status,
		Length:  byte(len(data)),
		Data:    data,
	}
	if _, err := writer.Write(header.Bytes()); err != nil {
		return err
	}
	return nil
}
func readContent(stream io.ReadWriteCloser, buf chan []byte) {
	defer stream.Close()
	for {
		var bf [readLen]byte
		n, err := stream.Read(bf[:])
		if err != nil {
			log.Trace("<---socksv server:", err)
			return
		}
		buf <- bf[0:n]
	}
}

type RelayStream struct {
	Addr string
	req  *socks5.Request
	//connection with socks5 client
	conn *net.TCPConn
}

func NewRelayStream(addr string, req *socks5.Request, inConn *net.TCPConn) *RelayStream {
	return &RelayStream{
		Addr: addr,
		req:  req,
		conn: inConn,
	}
}
func NewRelayStreamServer() *RelayStream {
	return &RelayStream{
		Addr: "",
	}
}
func (s *RelayStream) ID() protocol.ProtocolId {
	return protocol.Relay
}
func (s *RelayStream) In(rw io.ReadWriteCloser, session *smux.Session) error {
	//write reply to socks5 client
	if err := s.writeToSocks5(session); err != nil {
		return err
	}
	//write request to proxy server
	if err := writeRequest(rw, CmdConnect, []byte(s.Addr)); err != nil {
		return err
	}
	var rep Reply
	if err := rep.Read(rw); err != nil {
		return err
	}
	if rep.Status != StatusSuccess {
		return errors.New("connect failed")
	}
	log.Info("accept: ", s.Addr)
	//exchange data:socks5 client<--->proxy server
	//defer rw.Close()
	//defer s.conn.Close()
	protocol.ExchangeData(s.conn, rw)
	return nil
}
func (s *RelayStream) writeToSocks5(session *smux.Session) error {
	a, addr, port, err := socks5.ParseAddress(session.LocalAddr().String())
	if err != nil {
		socks5.ReplyError(s.req, s.conn, socks5.RepHostUnreachable)
		return err
	}
	successRep := socks5.NewReply(socks5.RepSuccess, a, addr, port)
	buf := bytes.NewBuffer(nil)
	if err := successRep.Write(buf); err != nil {
		return err
	}
	if _, err := s.conn.Write(buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (s *RelayStream) Out(rw io.ReadWriteCloser, session *smux.Session) error {
	defer rw.Close()
	targetAddr, err := acceptConnect(rw)
	if err != nil {
		log.Warn(err)
		return err
	}
	tunnel(rw, targetAddr)
	return nil
}
func tunnel(rw io.ReadWriteCloser, targetAddr string) {
	tmp, err := net.Dial("tcp", targetAddr)
	if err != nil {
		_ = writeReply(rw, StatusTargetUnreachable, nil)
		log.Warn("dial "+targetAddr+" error:", err)
		return
	}
	log.Info("dial:", targetAddr)
	conn := tmp.(*net.TCPConn)

	_ = writeReply(rw, StatusSuccess, nil)
	//defer conn.Close()
	//defer rw.Close()
	//exchange data:proxy server<--->target website
	protocol.ExchangeData(rw, conn)
}

func acceptConnect(stream io.ReadWriteCloser) (string, error) {
	var req Request
	if err := req.Read(stream); err != nil {
		return "", err
	}
	if req.Version != Version1 {
		_ = writeReply(stream, StatusUnSupportedVersion, nil)
		return "", errors.New("StatusUnSupportedVersion")
	}
	if req.Cmd != CmdConnect {
		_ = writeReply(stream, StatusUnsupportedCmd, nil)
		return "", errors.New("StatusUnsupportedCmd")
	}
	if req.Encrypt != EncryptAES {
		_ = writeReply(stream, StatusUnSupportedEncrypt, nil)
		return "", errors.New("StatusUnSupportedEncrypt")
	}
	targetAddr := string(req.Data)
	return targetAddr, nil
}
