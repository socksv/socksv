package protocol

import (
	"io"
)
import log "github.com/sirupsen/logrus"

///client <---> socks server <---> target
func ExchangeData(client io.ReadWriter, target io.ReadWriter) {
	go func() {
		var bf [1024 * 2]byte
		for {
			//read from target server
			n, err := target.Read(bf[:])
			if err != nil {
				log.Trace("<---target:", err)
				return
			}
			//write to client
			if _, err := client.Write(bf[0:n]); err != nil {
				log.Trace("--->client:", err)
				return
			}
		}
	}()
	var bf [1024 * 2]byte
	for {
		//read the request from client and send it to target server
		i, err := client.Read(bf[:])
		if err != nil {
			log.Trace("<---client:", err)
			return
		}
		if _, err := target.Write(bf[0:i]); err != nil {
			log.Trace("--->server:", err)
			return
		}
	}
}
