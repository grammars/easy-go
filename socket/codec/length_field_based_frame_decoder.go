package socket

import (
	"encoding/binary"
	"io"
)

type LengthFieldBasedFrameDecoder struct {
	ByteOrder         binary.ByteOrder
	MaxFrameLength    int
	LengthFieldOffset int
	LengthFieldLength int
	LengthAdjustment  int
	//InitialBytesToStrip int 不需要该字段，因为此处设计将以length为中心 切割成 2部分
}

func (decoder *LengthFieldBasedFrameDecoder) Decode(reader io.Reader) ([]byte, []byte, error) {
	var headerBytes []byte = nil
	if decoder.LengthFieldOffset > 0 {
		headerBytes = make([]byte, decoder.LengthFieldOffset)
		if _, err := io.ReadFull(reader, headerBytes); err != nil {
			return nil, nil, err
		}
	}
	lengthBuffer := make([]byte, decoder.LengthFieldLength)
	if _, err := io.ReadFull(reader, lengthBuffer); err != nil {
		return nil, nil, err
	}

	// length 之后的 内容长度
	contentLength := int(decoder.ByteOrder.Uint32(lengthBuffer)) + decoder.LengthAdjustment

	contentBytes := make([]byte, contentLength)
	if _, err := io.ReadFull(reader, contentBytes); err != nil {
		return headerBytes, nil, err
	}

	return headerBytes, contentBytes, nil
}
