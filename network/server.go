package network

import (
	log "github.com/sirupsen/logrus"
	"net"
	"socksv/network/smux"
	"socksv/protocol"
)

type Server struct {
	addr     *net.TCPAddr
	session  *smux.Session
	streams  map[protocol.ProtocolID]*smux.Stream
	handlers map[protocol.ProtocolID]protocol.Protocol
}

func NewServer(addr string) (*Server, error) {
	taddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &Server{
		addr:     taddr,
		session:  nil,
		streams:  make(map[protocol.ProtocolID]*smux.Stream),
		handlers: make(map[protocol.ProtocolID]protocol.Protocol),
	}, nil
}
func (s *Server) AddStreamHandler(handler protocol.Protocol) {
	s.handlers[handler.ID()] = handler
}
func (s *Server) Listen() {

	listener, _ := net.ListenTCP("tcp", s.addr)
	defer listener.Close()
	log.Info("server  started on ", s.addr.String(), " ...")
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Warn(err)
			continue
		}
		log.Debug("client connected: " + conn.RemoteAddr().String())
		go s.handle(conn)
	}
}

func (s *Server) handle(conn *net.TCPConn) {
	session, err := smux.Server(conn, nil)
	if err != nil {
		log.Warn(err)
		return
	}
	//TODO:how to close connection,Session??
	defer session.Close()
	defer conn.Close()
	for {
		stream, err := session.AcceptStream()
		if err != nil {
			log.Warn(err)
			break
		}
		handler, ok := s.handlers[stream.ProtocolID()]
		if !ok {
			log.Warn("protocol not supported:", stream.ID())
			break
		}
		go handler.Out(stream, session)
	}
}
