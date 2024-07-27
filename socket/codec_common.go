package socket

import (
	"io"
)

type CodecResult struct {
	Overflow    bool
	FrameLength int
	HeaderBytes []byte
	BodyBytes   []byte
}

type FrameDecoder[VD any] interface {
	Decode(visitor *Visitor[VD], reader io.Reader) (CodecResult, error)
}
