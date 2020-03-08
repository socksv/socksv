package network

import (
	log "github.com/sirupsen/logrus"
	"net"
	"socksv/network/smux"
	"socksv/protocol"
)

type Client struct {
	serverAddr *net.TCPAddr
	session    *smux.Session
}

func NewClient(addr string) (*Client, error) {
	taddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	c := &Client{
		serverAddr: taddr,
		session:    nil,
	}
	err = c.Connect()
	if err != nil {
		return nil, err
	}
	return c, nil
}
func (c *Client) Connect() error {

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
	c.session = session
	return err
}
func (c *Client) Open(handler protocol.StreamHandler) error {
	stream, err := c.session.OpenStreamFix(handler.ID())
	if err != nil {
		return err
	}
	if err := handler.In(stream); err != nil {
		return err
	}
	return nil
}
