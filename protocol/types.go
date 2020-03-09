package protocol

import (
	"io"
	"socksv/network/smux"
)

type ProtocolId = uint32
type StreamHandler interface {
	ID() ProtocolId
	//client side
	//remember to close rw
	In(rw io.ReadWriteCloser, session *smux.Session) error
	//server side
	//remember to close rw
	Out(rw io.ReadWriteCloser, session *smux.Session) error
}

type Protocol struct {
	smux.Stream
	pid uint32
}
