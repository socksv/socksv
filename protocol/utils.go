package protocol

import (
	"io"
)
import log "github.com/sirupsen/logrus"

///client <---> middle <---> target
func ExchangeData(client io.ReadWriteCloser, target io.ReadWriteCloser) {
	defer client.Close()
	defer target.Close()
	iseof := false
	go func() {
		var bf [1024 * 2]byte
		for {
			//read from target server
			n, err := target.Read(bf[:])
			if err != nil {
				if !iseof {
					log.Trace("<---target:", err)
				}
				if err == io.EOF {
					iseof = true
				}
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
			if !iseof {
				log.Trace("<---client:", err)
			}
			if err == io.EOF {
				iseof = true
			}
			return
		}
		if _, err := target.Write(bf[0:i]); err != nil {
			log.Trace("--->server:", err)
			return
		}
	}
}
