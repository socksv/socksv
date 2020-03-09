package protocol

import (
	"io"
	"socksv/network/smux"
)

type ProtocolID = byte
type Protocol interface {
	ID() ProtocolID
	//client side
	//remember to close rw
	In(rw io.ReadWriteCloser, session *smux.Session) error
	//server side
	//remember to close rw
	Out(rw io.ReadWriteCloser, session *smux.Session) error
}
