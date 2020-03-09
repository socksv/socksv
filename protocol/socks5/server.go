package socks5

import (
	"bytes"
	log "github.com/sirupsen/logrus"
	"net"
)

var SupportedCommands = []byte{CmdConnect, CmdUDP}
var ConnectHandler = DirectConnect

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
	log.Info("socks5 server started on ", s.TCPAddr.String(), " ...")
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
		ReplyError(&req, conn, RepCommandNotSupported)
		conn.Close()
	} else {
		if req.Cmd == CmdConnect {
			defer conn.Close()
			err := ConnectHandler(&req, conn)
			if err != nil {
				log.Warn("connect target server error:", err)
				return
			}
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
