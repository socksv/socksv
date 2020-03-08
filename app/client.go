package app

import (
	"github.com/sirupsen/logrus"
	"net"
	"socksv/network"
	"socksv/protocol/relay"
	"socksv/protocol/socks5"
)

type Client struct {
	socks5Server *socks5.Server
	proxyClient  *network.Client
}

func NewClient(socks5Addr, proxyAddr string) *Client {
	server, err := socks5.NewServer(socks5Addr, "", "", 100, 100)
	if err != nil {
		panic(err)
	}
	client, err := network.NewClient(proxyAddr)
	if err != nil {
		panic(err)
	}
	c := &Client{
		socks5Server: server,
		proxyClient:  client,
	}
	socks5.ConnectHandler = c.ProxyConnect
	return c
}
func (c *Client) ProxyConnect(req *socks5.Request, inConn *net.TCPConn) error {
	stream := relay.NewRelayStream(req.Address(), req, inConn)
	err := c.proxyClient.Open(stream)
	if err != nil {
		socks5.ReplyError(req, inConn, socks5.RepHostUnreachable)
		logrus.Warn(err)
	}
	return nil
}

//accept socks5 inbound stream
func (c *Client) Accept() {
	c.socks5Server.Listen()
}
