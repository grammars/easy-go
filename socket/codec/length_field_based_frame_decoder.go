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

func (decoder *LengthFieldBasedFrameDecoder) Decode(reader io.Reader) (*CodecResult, error) {
	result := &CodecResult{}
	if decoder.LengthFieldOffset > 0 {
		result.HeaderBytes = make([]byte, decoder.LengthFieldOffset)
		if _, err := io.ReadFull(reader, result.HeaderBytes); err != nil {
			return result, err
		}
	}
	lengthBuffer := make([]byte, decoder.LengthFieldLength)
	if _, err := io.ReadFull(reader, lengthBuffer); err != nil {
		return result, err
	}

	// length 之后的 内容长度
	bodyLength := int(decoder.ByteOrder.Uint32(lengthBuffer)) + decoder.LengthAdjustment

	result.BodyBytes = make([]byte, bodyLength)
	if _, err := io.ReadFull(reader, result.BodyBytes); err != nil {
		return result, err
	}

	result.FrameLength = len(result.HeaderBytes) + decoder.LengthFieldLength + len(result.BodyBytes)

	return result, nil
}
