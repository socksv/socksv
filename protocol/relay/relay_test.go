package relay

import (
	"github.com/sirupsen/logrus"
	"net"
	"socksv/network"
	"socksv/protocol/ping"
	"testing"
)

var rs = NewRelayStream("180.101.49.12:443")
var pi = ping.NewPing()

func TestServer(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	server, err := network.NewServer("0.0.0.0:8888")
	if err != nil {
		panic(err)
	}
	server.AddStreamHandler(pi)
	server.AddStreamHandler(rs)
	server.Listen()
}
func TestClient(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
	client, err := network.NewClient("127.0.0.1:8888")
	if err != nil {
		panic(err)
	}
	go client.Open(pi)
	err = client.Open(rs)
	if err != nil {
		panic(err)
	}
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
