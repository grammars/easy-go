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

// Combine 把 HeaderBytes 和 BodyBytes 合并起来
func (cr *CodecResult) Combine() []byte {
	return append(cr.HeaderBytes, cr.BodyBytes...)
}

type FrameDecoder interface {
	Decode(visitor VisitorSupport, reader io.Reader) (CodecResult, error)
}
