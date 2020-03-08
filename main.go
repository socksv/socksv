package main

import (
	"flag"
	"fmt"
	"github.com/sirupsen/logrus"
	"socksv/protocol/relay"
	"socksv/protocol/socks5"
)

var logLevel = logrus.DebugLevel
var port string
var lv int
var server string
var target string

func main() {
	flag.StringVar(&port, "p", "1080", `server port.`)
	flag.IntVar(&lv, "l", 1, `log level.0-info;1-debug;2-trace;3-warn;4-error.`)
	flag.StringVar(&server, "s", "", "relay server to connect.")
	flag.StringVar(&target, "x", "", "target address to relay.")
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
		//startRelayServer("0.0.0.0:" + port)
		startSocks5Server("0.0.0.0:" + port)
	} else if target != "" {
		client, err := relay.NewClient(server)
		if err != nil {
			panic(err)
		}
		recv, err := client.Open(server)
		if err != nil {
			panic(err)
		}
		for {
			data := <-recv
			logrus.Infof("rcv:%x\n", data)
		}

	} else {
		fmt.Println("please run as server or client")
	}

}
func startSocks5Server(addr string) {
	server, err := socks5.NewServer(addr, "", "", 100, 100)
	if err != nil {
		panic(err)
	}
	server.Listen()
}
func startRelayServer(addr string) {
	server, err := relay.NewServer(addr)
	if err != nil {
		panic(err)
	}
	server.Listen()
}
