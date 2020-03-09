package sv

import (
	"github.com/sirupsen/logrus"
	"net"
	"socksv/network"
	"socksv/protocol/ping"
	"testing"
)

var pi = ping.NewPing()

func TestServer(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	server, err := network.NewServer("0.0.0.0:8080")
	if err != nil {
		panic(err)
	}
	//server.AddStreamHandler(pi)
	server.AddStreamHandler(NewSocksVProtocolEmpty())
	server.Listen()
}
func TestDial(t *testing.T) {
	addr := "https://www.baidu.com"
	tmp, err := net.Dial("tcp", addr)
	if err != nil {
		panic(err)
	}
	//addrs,err:=net.LookupHost(addr)
	//if err != nil {
	//	panic(err)
	//}
	println("dial:", tmp)
}
