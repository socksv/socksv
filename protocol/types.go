package protocol

import (
	"io"
	"socksv/network/smux"
)

type ProtocolId = uint32
type StreamHandler interface {
	ID() ProtocolId
	//client side
	In(rw io.ReadWriteCloser, session *smux.Session) error
	//server side
	Out(rw io.ReadWriteCloser, session *smux.Session) error
}
