package app

import (
	"github.com/sirupsen/logrus"
	"net"
	"socksv/network"
	"socksv/protocol/ping"
	"socksv/protocol/sv"
	"socksv/socks5"
)

var EnablePing = false

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
	if EnablePing {
		go client.Open(ping.NewPing())
	}
	c := &Client{
		socks5Server: server,
		proxyClient:  client,
	}
	//set socks5 handler to ProxyConnect
	socks5.ConnectHandler = c.proxyConnect
	return c
}

//proxyConnect is used in socks5
func (c *Client) proxyConnect(req *socks5.Request, inConn *net.TCPConn) error {
	stream := sv.NewSocksVProtocol(req.Address(), req, inConn)
	err := c.proxyClient.Open(stream)
	if err != nil {
		socks5.ReplyError(req, inConn, socks5.RepHostUnreachable)
		logrus.Warn(req.Address()+": ", err)
	}
	return nil
}

//accept socks5 inbound stream
func (c *Client) Accept() {
	c.socks5Server.Listen()
}
func StartProxyClient(socksListenAddr, proxyServerAddr string) {
	client := NewClient(socksListenAddr, proxyServerAddr)
	client.Accept()
}
func StartProxyServer(addr string) {
	server, err := network.NewServer(addr)
	if err != nil {
		panic(err)
	}
	if EnablePing {
		go server.AddStreamHandler(ping.NewPing())
	}
	server.AddStreamHandler(sv.NewSocksVProtocolEmpty())
	server.Listen()
}
