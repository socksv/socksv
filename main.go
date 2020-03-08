package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"socksv/app"
	"socksv/network"
	"socksv/protocol/relay"
	"socksv/protocol/socks5"
)

var logLevel = logrus.DebugLevel
var port string
var serverPort string
var lv int
var server string

//var target string

func main() {
	flag.StringVar(&port, "p", "1080", `socks5 server port.`)
	flag.StringVar(&serverPort, "P", "8080", `proxy server port.`)
	flag.IntVar(&lv, "l", 1, `log level.0-info;1-debug;2-trace;3-warn;4-error.`)
	flag.StringVar(&server, "s", "", "relay server to connect.")
	//flag.StringVar(&target, "x", "", "target address to relay.")
	flag.Parse()
	switch lv {
	case 0:
		logLevel = logrus.InfoLevel
	case 1:
		logLevel = logrus.DebugLevel
	case 2:
		logLevel = logrus.TraceLevel
	case 3:
		logLevel = logrus.WarnLevel
	case 4:
		logLevel = logrus.ErrorLevel
	}
	logrus.SetLevel(logLevel)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})
	if server == "" {
		StartProxyServer("0.0.0.0:" + serverPort)
	} else {
		//go StartSocks5Server("0.0.0.0:" + port)
		//client := StartProxyClient(server)
		//err := client.Open(relay.NewRelayStream("180.101.49.12:443"))
		//if err != nil {
		//	panic(err)
		//}
		client := app.NewClient("0.0.0.0:"+port, server)
		client.Accept()
	}

}
func StartSocks5Server(addr string) {
	server, err := socks5.NewServer(addr, "", "", 100, 100)
	if err != nil {
		panic(err)
	}
	server.Listen()
}
func StartProxyServer(addr string) {
	server, err := network.NewServer(addr)
	if err != nil {
		panic(err)
	}
	server.AddStreamHandler(relay.NewRelayStreamServer())
	server.Listen()
}
func StartProxyClient(addr string) *network.Client {
	client, err := network.NewClient(addr)
	if err != nil {
		panic(err)
	}
	return client
}
