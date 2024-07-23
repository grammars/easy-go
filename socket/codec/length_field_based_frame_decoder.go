package socket

import (
	"encoding/binary"
	"fmt"
	"io"
	"log/slog"
)

type LengthFieldBasedFrameDecoder struct {
	ByteOrder         binary.ByteOrder
	MaxFrameLength    int
	LengthFieldOffset int
	LengthFieldLength int
	LengthAdjustment  int
	//InitialBytesToStrip int 不需要该字段，因为此处设计将以length为中心 切割成 2部分
}

func (decoder *LengthFieldBasedFrameDecoder) Decode(reader io.Reader) (CodecResult, error) {
	slog.Info("准备解码")
	result := CodecResult{}
	if decoder.LengthFieldOffset > 0 {
		result.HeaderBytes = make([]byte, decoder.LengthFieldOffset)
		if _, err := io.ReadFull(reader, result.HeaderBytes); err != nil {
			return result, err
		}
	}
	slog.Info("读取到HeaderBytes", "HeaderBytes", fmt.Sprintf("%x", result.HeaderBytes))

	lengthBuffer := make([]byte, decoder.LengthFieldLength)
	if _, err := io.ReadFull(reader, lengthBuffer); err != nil {
		return result, err
	}

	slog.Info("读取到lengthBuffer")
	// length 之后的 内容长度
	bodyLength := int(decoder.ByteOrder.Uint32(lengthBuffer)) + decoder.LengthAdjustment
	slog.Info("读取到bodyLength", "bodyLength", bodyLength)

	calcFrameLength := bodyLength + decoder.LengthFieldOffset + decoder.LengthFieldLength
	if calcFrameLength > decoder.MaxFrameLength {
		slog.Error("数据帧溢出", "预计帧长度", calcFrameLength, "最大允许帧长度", decoder.MaxFrameLength)
		result.Overflow = true
		return result, nil
	}

	result.BodyBytes = make([]byte, bodyLength)
	if _, err := io.ReadFull(reader, result.BodyBytes); err != nil {
		return result, err
	}

	slog.Info("读取到BodyBytes", "BodyBytes", fmt.Sprintf("%x", result.BodyBytes))

	result.FrameLength = len(result.HeaderBytes) + decoder.LengthFieldLength + len(result.BodyBytes)

	return result, nil
}
