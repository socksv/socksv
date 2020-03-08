package ping

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"io"
	"socksv/protocol"
	"time"
)

const (
	ping = 1
	pong = 2
)

var Timeout int64 = 30

type Ping struct {
	timeout time.Duration
}

func NewPing() *Ping {
	return &Ping{
		timeout: time.Duration(Timeout),
	}
}
func (p *Ping) ID() protocol.ProtocolId {
	return protocol.Ping
}
func (p *Ping) In(rw io.ReadWriteCloser) error {
	defer rw.Close()
	ticker := time.Tick(1 * time.Second)
	timeout := time.After(p.timeout * time.Second)
	go func() {
		for {
			po := make([]byte, 1)
			if _, err := rw.Read(po); err != nil {
				log.Warn("read error:", err)
				return
			}
			if po[0] != pong {
				log.Warn("unknown command:", po[0])
				return
			}
			log.Trace("<---pong")
			//reset timeout
			timeout = time.After(30 * time.Second)
		}
	}()
	for {
		select {
		case <-ticker:
			if _, err := rw.Write([]byte{ping}); err != nil {
				log.Warn("write error:", err)
				return err
			}
			log.Trace("--->ping")
		case <-timeout:
			log.Warn("receive from server timeout")
			return errors.New("receive timeout")
		}
	}

	return nil
}
func (p *Ping) Out(rw io.ReadWriteCloser) error {
	defer rw.Close()
	for {
		pi := make([]byte, 1)
		if _, err := rw.Read(pi); err != nil {
			log.Warn("read error:", err)
			return err
		}
		if pi[0] == ping {
			log.Trace("<---ping")
			if _, err := rw.Write([]byte{pong}); err != nil {
				log.Warn("write error:", err)
				return err
			}
			log.Trace("--->pong")
		} else {
			log.Warn("unknown command:", pi[0])
			return errors.New("unknown command")
		}
	}
	return nil
}
