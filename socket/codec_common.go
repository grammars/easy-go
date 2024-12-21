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

type FrameDecoder interface {
	Decode(visitor VisitorSupport, reader io.Reader) (CodecResult, error)
}
