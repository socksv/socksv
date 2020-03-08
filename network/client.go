package network

import (
	log "github.com/sirupsen/logrus"
	"net"
	"socksv/network/smux"
	"socksv/protocol"
)

type Client struct {
	serverAddr *net.TCPAddr
	Session    *smux.Session
}

func NewClient(addr string) (*Client, error) {
	taddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	c := &Client{
		serverAddr: taddr,
		Session:    nil,
	}
	err = c.connect()
	if err != nil {
		return nil, err
	}
	return c, nil
}
func (c *Client) connect() error {

	conn, err := net.Dial("tcp", c.serverAddr.String())
	if err != nil {
		return err
	}
	log.Info("dial success:", c.serverAddr.String())
	// Setup client side of smux
	session, err := smux.Client(conn, nil)
	if err != nil {
		return err
	}
	c.Session = session
	return err
}
func (c *Client) Open(handler protocol.StreamHandler) error {
	stream, err := c.Session.OpenStreamFix(handler.ID())
	if err != nil {
		return err
	}
	if err := handler.In(stream, c.Session); err != nil {
		return err
	}
	return nil
}
