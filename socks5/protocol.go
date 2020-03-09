package socks5

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/sirupsen/logrus"
	"io"
	"net"
	"socksv/utils"
	"strconv"
)

const (
	// Ver is socks protocol version
	Ver byte = 0x05

	// MethodNone is none method
	MethodNone byte = 0x00
	// MethodGSSAPI is gssapi method
	MethodGSSAPI byte = 0x01 // MUST support // TODO
	// MethodUsernamePassword is username/assword auth method
	MethodUsernamePassword byte = 0x02 // SHOULD support
	// MethodUnsupportAll means unsupport all given methods
	MethodUnsupportAll byte = 0xFF

	// UserPassVer is username/password auth protocol version
	UserPassVer byte = 0x01
	// UserPassStatusSuccess is success status of username/password auth
	UserPassStatusSuccess byte = 0x00
	// UserPassStatusFailure is failure status of username/password auth
	UserPassStatusFailure byte = 0x01 // just other than 0x00

	// CmdConnect is connect command
	CmdConnect byte = 0x01
	// CmdBind is bind command
	CmdBind byte = 0x02
	// CmdUDP is UDP command
	CmdUDP byte = 0x03

	// ATYPIPv4 is ipv4 address type
	ATYPIPv4 byte = 0x01 // 4 octets
	// ATYPDomain is domain address type
	ATYPDomain byte = 0x03 // The first octet of the address field contains the number of octets of name that follow, there is no terminating NUL octet.
	// ATYPIPv6 is ipv6 address type
	ATYPIPv6 byte = 0x04 // 16 octets

	// RepSuccess means that success for repling
	RepSuccess byte = 0x00
	// RepServerFailure means the server failure
	RepServerFailure byte = 0x01
	// RepNotAllowed means the request not allowed
	RepNotAllowed byte = 0x02
	// RepNetworkUnreachable means the network unreachable
	RepNetworkUnreachable byte = 0x03
	// RepHostUnreachable means the host unreachable
	RepHostUnreachable byte = 0x04
	// RepConnectionRefused means the connection refused
	RepConnectionRefused byte = 0x05
	// RepTTLExpired means the TTL expired
	RepTTLExpired byte = 0x06
	// RepCommandNotSupported means the request command not supported
	RepCommandNotSupported byte = 0x07
	// RepAddressNotSupported means the request address not supported
	RepAddressNotSupported byte = 0x08
)

// NegotiationRequest is the negotiation reqeust packet
type NegotiationRequest struct {
	Ver      byte
	NMethods byte
	Methods  []byte // 1-255 bytes
}

func (p *NegotiationRequest) read(writer io.Writer) error {
	w := bufio.NewWriter(writer)
	if err := w.WriteByte(p.Ver); err != nil {
		return err
	}
	if err := w.WriteByte(p.NMethods); err != nil {
		return err
	}
	if _, err := w.Write(p.Methods); err != nil {
		return err
	}
	return nil
}
func (p *NegotiationRequest) write(reader io.Reader) error {
	bb := make([]byte, 2)
	if _, err := io.ReadFull(reader, bb); err != nil {
		return err
	}
	p.Ver = bb[0]
	p.NMethods = bb[1]
	by := make([]byte, p.NMethods)
	if _, err := io.ReadFull(reader, by); err != nil {
		return err
	} else {
		p.Methods = by
	}
	return nil
}

// NegotiationReply is the negotiation reply packet
type NegotiationReply struct {
	Ver    byte
	Method byte
}

func (p *NegotiationReply) write(writer io.Writer) error {
	if _, err := writer.Write([]byte{p.Ver, p.Method}); err != nil {
		return err
	}
	return nil
}
func (p *NegotiationReply) read(reader io.Reader) error {
	r := bufio.NewReader(reader)
	if by, err := r.ReadByte(); err != nil {
		return err
	} else {
		p.Ver = by
	}
	if by, err := r.ReadByte(); err != nil {
		return err
	} else {
		p.Method = by
	}
	return nil
}

// UserPassNegotiationRequest is the negotiation username/password reqeust packet
type UserPassNegotiationRequest struct {
	Ver    byte
	Ulen   byte
	Uname  []byte // 1-255 bytes
	Plen   byte
	Passwd []byte // 1-255 bytes
}

func (p *UserPassNegotiationRequest) write(writer io.Writer) error {
	w := bufio.NewWriter(writer)
	if err := w.WriteByte(p.Ver); err != nil {
		return err
	}
	if err := w.WriteByte(p.Ulen); err != nil {
		return err
	}
	if _, err := w.Write(p.Uname); err != nil {
		return err
	}
	if err := w.WriteByte(p.Plen); err != nil {
		return err
	}
	if _, err := w.Write(p.Passwd); err != nil {
		return err
	}
	return nil
}
func (p *UserPassNegotiationRequest) read(reader io.Reader) error {
	r := bufio.NewReader(reader)
	if by, err := r.ReadByte(); err != nil {
		return err
	} else {
		p.Ver = by
	}
	if by, err := r.ReadByte(); err != nil {
		return err
	} else {
		p.Ulen = by
	}
	if by, err := r.ReadBytes(p.Ulen); err != nil {
		return err
	} else {
		p.Uname = by
	}
	if by, err := r.ReadByte(); err != nil {
		return err
	} else {
		p.Plen = by
	}
	if by, err := r.ReadBytes(p.Plen); err != nil {
		return err
	} else {
		p.Passwd = by
	}
	return nil
}

// UserPassNegotiationReply is the negotiation username/password reply packet
type UserPassNegotiationReply struct {
	Ver    byte
	Status byte
}

func (p *UserPassNegotiationReply) write(writer io.Writer) error {
	w := bufio.NewWriter(writer)
	if err := w.WriteByte(p.Ver); err != nil {
		return err
	}
	if err := w.WriteByte(p.Status); err != nil {
		return err
	}
	return nil
}
func (p *UserPassNegotiationReply) read(reader io.Reader) error {
	r := bufio.NewReader(reader)
	if by, err := r.ReadByte(); err != nil {
		return err
	} else {
		p.Ver = by
	}
	if by, err := r.ReadByte(); err != nil {
		return err
	} else {
		p.Status = by
	}
	return nil
}

// Request is the request packet
type Request struct {
	Ver     byte
	Cmd     byte
	Rsv     byte // 0x00
	Atyp    byte
	DstAddr []byte
	DstPort []byte // 2 bytes
}

// Address return request address like ip:xx
func (r *Request) Address() string {
	var s string
	if r.Atyp == ATYPDomain {
		s = bytes.NewBuffer(r.DstAddr[1:]).String()
	} else {
		s = net.IP(r.DstAddr).String()
	}
	p := strconv.Itoa(int(binary.BigEndian.Uint16(r.DstPort)))
	return net.JoinHostPort(s, p)
}

func (p *Request) write(writer io.Writer) error {
	w := bufio.NewWriter(writer)
	if err := w.WriteByte(p.Ver); err != nil {
		return err
	}
	if err := w.WriteByte(p.Cmd); err != nil {
		return err
	}
	if err := w.WriteByte(p.Rsv); err != nil {
		return err
	}
	if err := w.WriteByte(p.Atyp); err != nil {
		return err
	}
	if _, err := w.Write(p.DstAddr); err != nil {
		return err
	}
	if _, err := w.Write(p.DstPort); err != nil {
		return err
	}
	return nil
}
func (p *Request) read(r io.Reader) error {
	bb := make([]byte, 4)
	if _, err := io.ReadFull(r, bb); err != nil {
		return nil
	}
	p.Ver = bb[0]
	p.Cmd = bb[1]
	p.Rsv = bb[2]
	p.Atyp = bb[3]
	var n byte
	switch p.Atyp {
	case ATYPIPv4: //ipv4
		n = 4
		break
	case ATYPDomain: //domain name
		nb := make([]byte, 1)
		if _, err := io.ReadFull(r, nb); err != nil {
			return err
		} else {
			n = nb[0]
		}
		break
	case ATYPIPv6: //ipv6
		n = 16
		break
	default:
		return errors.New("unknow address type")

	}
	addr := make([]byte, n)
	if _, err := io.ReadFull(r, addr); err != nil {
		return err
	}
	if p.Atyp == ATYPDomain {
		p.DstAddr = append([]byte{n}, addr...)
	} else {
		p.DstAddr = addr
	}

	port := make([]byte, 2)
	if _, err := io.ReadFull(r, port); err != nil {
		return err
	}
	p.DstPort = port
	return nil
}

// Reply is the reply packet
type Reply struct {
	Ver  byte
	Rep  byte
	Rsv  byte // 0x00
	Atyp byte
	// CONNECT socks server's address which used to connect to dst addr
	// BIND ...
	// UDP socks server's address which used to connect to dst addr
	BndAddr []byte
	// CONNECT socks server's port which used to connect to dst addr
	// BIND ...
	// UDP socks server's port which used to connect to dst addr
	BndPort []byte // 2 bytes
}

// NewReply return reply packet can be writed into client
func NewReply(rep byte, atyp byte, bndaddr []byte, bndport []byte) *Reply {
	if atyp == ATYPDomain {
		bndaddr = append([]byte{byte(len(bndaddr))}, bndaddr...)
	}
	return &Reply{
		Ver:     Ver,
		Rep:     rep,
		Rsv:     0x00,
		Atyp:    atyp,
		BndAddr: bndaddr,
		BndPort: bndport,
	}
}

func (p *Reply) Write(w io.Writer) error {
	if _, err := w.Write([]byte{p.Ver, p.Rep, p.Rsv, p.Atyp}); err != nil {
		return err
	}
	if _, err := w.Write(p.BndAddr); err != nil {
		return err
	}
	if _, err := w.Write(p.BndPort); err != nil {
		return err
	}
	return nil
}
func (p *Reply) Read(reader io.Reader) error {
	r := bufio.NewReader(reader)
	if by, err := r.ReadByte(); err != nil {
		return err
	} else {
		p.Ver = by
	}
	if by, err := r.ReadByte(); err != nil {
		return err
	} else {
		p.Rep = by
	}
	if by, err := r.ReadByte(); err != nil {
		return err
	} else {
		p.Rsv = by
	}
	if by, err := r.ReadByte(); err != nil {
		return err
	} else {
		p.Atyp = by
	}
	var n byte
	switch p.Atyp {
	case ATYPIPv4: //ipv4
		n = 4
		break
	case ATYPDomain: //domain name
		if by, err := r.ReadByte(); err != nil {
			return err
		} else {
			n = by
		}
		break
	case ATYPIPv6: //ipv6
		n = 16
		break
	default:
		return errors.New("unknow aty type")

	}
	if by, err := r.ReadBytes(n); err != nil {
		return err
	} else {
		p.BndAddr = append([]byte{n}, by...)
	}
	if by, err := r.ReadBytes(2); err != nil {
		return err
	} else {
		p.BndPort = by
	}
	return nil
}

// Datagram is the UDP packet
type Datagram struct {
	Rsv     []byte // 0x00 0x00
	Frag    byte
	Atyp    byte
	DstAddr []byte
	DstPort []byte // 2 bytes
	Data    []byte
}

func ReplyError(req *Request, inConn *net.TCPConn, cmd byte) {
	var rep *Reply
	if req.Atyp == ATYPIPv4 || req.Atyp == ATYPDomain {
		rep = NewReply(cmd, ATYPIPv4, []byte{0x00, 0x00, 0x00, 0x00}, []byte{0x00, 0x00})
	} else {
		rep = NewReply(cmd, ATYPIPv6, []byte(net.IPv6zero), []byte{0x00, 0x00})
	}
	buf := bytes.NewBuffer(nil)
	if err := rep.Write(buf); err != nil {
		logrus.Warn(err)
	}
	if _, err := inConn.Write(buf.Bytes()); err != nil {
		logrus.Warn(err)
	}
}

type Connect = func(req *Request, conn *net.TCPConn) error

func directConnect(req *Request, inConn *net.TCPConn) error {
	logrus.Info("Dial:", req.Address())
	tmp, err := net.Dial("tcp", req.Address())
	if err != nil {
		ReplyError(req, inConn, RepHostUnreachable)
		return err
	}
	outConn := tmp.(*net.TCPConn)
	a, addr, port, err := ParseAddress(outConn.LocalAddr().String())
	if err != nil {
		ReplyError(req, inConn, RepHostUnreachable)
		return err
	}
	successRep := NewReply(RepSuccess, a, addr, port)
	buf := bytes.NewBuffer(nil)
	if err := successRep.Write(buf); err != nil {
		return err
	}
	if _, err := inConn.Write(buf.Bytes()); err != nil {
		return err
	}
	utils.ProxyData(inConn, outConn)
	return nil
}
