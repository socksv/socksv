package main

import (
	"flag"
	"github.com/sirupsen/logrus"
	"path"
	"runtime"
	"socksv/app"
	"strconv"
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
	logrus.SetReportCaller(true)
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			//s := strings.Split(f.Function, ".")
			//funcname := s[len(s)-1]
			_, filename := path.Split(f.File)
			return "", "[" + filename + ":" + strconv.Itoa(f.Line) + "]"
		},
	})
	if server == "" {
		app.StartProxyServer("0.0.0.0:" + serverPort)
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
