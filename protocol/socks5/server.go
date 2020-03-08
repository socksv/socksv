package socks5

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"socksv/protocol"

	"net"
)

var SupportedCommands = []byte{CmdConnect, CmdUDP}

type Server struct {
	UserName    string
	Password    string
	Method      byte
	TCPDeadline int
	TCPTimeout  int
	TCPAddr     *net.TCPAddr
}

func NewServer(addr, uname, pwd string, tcpDeadline, tcpTimeout int) (*Server, error) {
	//if _, p, err := net.SplitHostPort(addr); err != nil {
	//	return nil, err
	//}
	taddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	m := MethodNone
	if uname != "" && pwd != "" {
		m = MethodUsernamePassword
	}
	s := Server{
		UserName:    uname,
		Password:    pwd,
		Method:      m,
		TCPDeadline: tcpDeadline,
		TCPTimeout:  tcpTimeout,
		TCPAddr:     taddr,
	}

	return &s, nil
}
func (s *Server) Listen() {

	tcpListener, _ := net.ListenTCP("tcp", s.TCPAddr)
	defer tcpListener.Close()
	log.Info("socks5 server  started on ", s.TCPAddr.String(), " ...")
	//listen
	for {
		conn, err := tcpListener.AcceptTCP()
		if err != nil {
			log.Warn(err)
			continue
		}
		log.Debug("client connected: " + conn.RemoteAddr().String())
		go handle(conn)
	}
}

func handle(conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()
	defer func() {
		log.Debug("client disconnected: " + ipStr)
		conn.Close()
	}()
	handleProtocol(conn)
}

func handleProtocol(conn *net.TCPConn) {
	err := negotiation(conn)
	if err != nil {
		log.Error("negotiation error:", err)
		return
	}
	request(conn)
}

func request(conn *net.TCPConn) {
	var req Request

	if err := req.read(conn); err != nil {
		log.Warn("request error:", err)
		return
	}
	log.Trace("connect request:", req)

	var supported bool
	for _, c := range SupportedCommands {
		if req.Cmd == c {
			supported = true
			break
		}
	}
	if !supported {
		replyError(&req, conn, RepCommandNotSupported)
	} else {
		if req.Cmd == CmdConnect {
			outConn, err := connect(&req, conn)
			if err != nil {
				log.Warn("connect target server error:", err)
				return
			}
			defer outConn.Close()
			defer conn.Close()
			protocol.ExchangeData(conn, outConn)
		}
	}

}
func negotiation(conn *net.TCPConn) error {
	var nreq NegotiationRequest
	if err := nreq.write(conn); err != nil {
		return err
	}
	log.Trace("negotiation request:", nreq)

	nrep := NegotiationReply{
		Ver:    Ver,
		Method: MethodNone,
	}
	buf := bytes.NewBuffer(nil)
	if err := nrep.write(buf); err != nil {
		return err
	}
	if _, err := conn.Write(buf.Bytes()); err != nil {
		return err
	}
	return nil
}
func connect(req *Request, inConn *net.TCPConn) (*net.TCPConn, error) {
	log.Info("Dial:", req.Address())
	tmp, err := net.Dial("tcp", req.Address())
	if err != nil {
		replyError(req, inConn, RepHostUnreachable)
		return nil, err
	}
	outConn := tmp.(*net.TCPConn)
	a, addr, port, err := parseAddress(outConn.LocalAddr().String())
	if err != nil {
		replyError(req, inConn, RepHostUnreachable)
		return nil, err
	}
	successRep := NewReply(RepSuccess, a, addr, port)
	buf := bytes.NewBuffer(nil)
	if err := successRep.write(buf); err != nil {
		return nil, err
	}
	if _, err := inConn.Write(buf.Bytes()); err != nil {
		return nil, err
	}
	return outConn, nil
}

func replyError(req *Request, inConn *net.TCPConn, cmd byte) {
	var rep *Reply
	if req.Atyp == ATYPIPv4 || req.Atyp == ATYPDomain {
		rep = NewReply(cmd, ATYPIPv4, []byte{0x00, 0x00, 0x00, 0x00}, []byte{0x00, 0x00})
	} else {
		rep = NewReply(cmd, ATYPIPv6, []byte(net.IPv6zero), []byte{0x00, 0x00})
	}
	buf := bytes.NewBuffer(nil)
	if err := rep.write(buf); err != nil {
		log.Warn(err)
	}
	if _, err := inConn.Write(buf.Bytes()); err != nil {
		log.Warn(err)
	}
}
