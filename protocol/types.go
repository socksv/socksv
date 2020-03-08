package protocol

import "io"

type ProtocolId = uint32
type StreamHandler interface {
	ID() ProtocolId
	//client side
	In(rw io.ReadWriteCloser) error
	//server side
	Out(rw io.ReadWriteCloser) error
}
