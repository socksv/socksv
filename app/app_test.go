package app

import (
	"github.com/sirupsen/logrus"
	"path"
	"runtime"
	"strconv"
	"testing"
)

func TestServer(t *testing.T) {
	logrus.SetLevel(logrus.InfoLevel)
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
	StartProxyServer("0.0.0.0:8080")
}
func TestClient(t *testing.T) {
	logrus.SetLevel(logrus.TraceLevel)
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
	client := NewClient("0.0.0.0:1080", "127.0.0.1:8080")
	client.Accept()
}
