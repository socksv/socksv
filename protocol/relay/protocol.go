package relay

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"io"
	"net"
	"socksv/protocol"
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
}

func NewRelayStream(addr string) *RelayStream {
	return &RelayStream{
		Addr: addr,
	}
}
func (s *RelayStream) ID() protocol.ProtocolId {
	return protocol.Relay
}
func (s *RelayStream) In(rw io.ReadWriteCloser) error {
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
	buf := make(chan []byte, 1000)
	readContent(rw, buf)
	return nil
}

func (s *RelayStream) Out(rw io.ReadWriteCloser) error {
	defer rw.Close()
	targetAddr, err := acceptConnect(rw)
	if err != nil {
		log.Warn(err)
		return err
	}
	tunnel(rw, targetAddr)
	return nil
}
func tunnel(stream io.ReadWriter, targetAddr string) {
	tmp, err := net.Dial("tcp", targetAddr)
	if err != nil {
		_ = writeReply(stream, StatusTargetUnreachable, nil)
		log.Warn("dial "+targetAddr+" error:", err)
		return
	}
	log.Info("dial:", targetAddr)
	conn := tmp.(*net.TCPConn)
	defer conn.Close()
	_ = writeReply(stream, StatusSuccess, nil)
	protocol.ExchangeData(stream, conn)
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
