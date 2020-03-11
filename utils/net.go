package utils

import (
	"encoding/binary"
	"github.com/sirupsen/logrus"
	"io"
)

//Proxy exchange data of client and target
//client <---> middle <---> target
//make sure both client and target reader and writer close,and return when
//stream is EOF
func ProxyData(client io.ReadWriteCloser, target io.ReadWriteCloser) {
	go func() {
		if err := WriteTo(target, client, false, false); err != nil {
			logrus.Trace("target ---> client error:", err)
		}
	}()
	err := WriteTo(client, target, false, false)
	if err != nil {
		logrus.Trace("client ---> target error:", err)
	}
}

//Read data from `from` and write to `to`
func WriteTo(from io.ReadWriter, to io.ReadWriter, rd, we bool) error {
	var bf [1024 * 2]byte
	for {
		var bb []byte
		if rd {
			if by, err := decrypt(from); err != nil {
				return err
			} else {
				bb = by
			}
		} else {
			i, err := from.Read(bf[:])
			if err != nil {
				return err
			}
			bb = bf[:i]
		}
		if we {
			//encrypt data
			bb = encrypt(bb)
		}
		if _, err := to.Write(bb); err != nil {
			return err
		}
	}
}

const lengthSize = 2

//decrypt with length
func decrypt(reader io.Reader) ([]byte, error) {
	byLen := make([]byte, lengthSize)
	if _, err := io.ReadFull(reader, byLen); err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint16(byLen)
	data := make([]byte, length)
	if _, err := io.ReadFull(reader, data); err != nil {
		return nil, err
	}
	//fmt.Printf("<---:%x\n", data)
	return AesDecrypt(data, Key), nil
}

//encrypt with length
func encrypt(data []byte) []byte {
	wr := AesEncrypt(data, Key)
	bb := make([]byte, len(wr)+lengthSize)
	binary.BigEndian.PutUint16(bb, uint16(len(wr)))
	copy(bb[lengthSize:], wr)
	//fmt.Printf("-->:%x\n", bb)
	return bb
}
