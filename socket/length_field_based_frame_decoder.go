package socket

import (
	"bytes"
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
	LengthAdjustment  int // 未严格验证
	//InitialBytesToStrip int 不需要该字段，因为此处设计将以length为中心 切割成 2部分
}

func (decoder *LengthFieldBasedFrameDecoder) Decode(visitor VisitorSupport, reader io.Reader) (CodecResult, error) {
	return decoder.Read(reader)
}

func (decoder *LengthFieldBasedFrameDecoder) ReadBytes(dataBytes []byte) (CodecResult, error) {
	return decoder.Read(bytes.NewBuffer(dataBytes))
}

func (decoder *LengthFieldBasedFrameDecoder) Read(reader io.Reader) (CodecResult, error) {
	if LogLevel <= 0 {
		slog.Info("LengthFieldBasedFrameDecoder准备解码")
	}
	result := CodecResult{}
	if decoder.LengthFieldOffset > 0 {
		result.HeaderBytes = make([]byte, decoder.LengthFieldOffset)
		if _, err := io.ReadFull(reader, result.HeaderBytes); err != nil {
			return result, err
		}
	}
	if LogLevel <= 0 {
		slog.Info("读取到HeaderBytes", "HeaderBytes", fmt.Sprintf("%x", result.HeaderBytes))
	}

	lengthBuffer := make([]byte, decoder.LengthFieldLength)
	if _, err := io.ReadFull(reader, lengthBuffer); err != nil {
		return result, err
	}
	var bl int
	if decoder.LengthFieldLength == 1 {
		bl = int(lengthBuffer[0])
	} else if decoder.LengthFieldLength == 2 {
		bl = int(decoder.ByteOrder.Uint16(lengthBuffer))
	} else if decoder.LengthFieldLength == 4 {
		bl = int(decoder.ByteOrder.Uint32(lengthBuffer))
	} else {
		return result, fmt.Errorf("不支持的LengthFieldLength=%d", decoder.LengthFieldLength)
	}

	if LogLevel <= 0 {
		slog.Info("读取到lengthBuffer")
	}

	// length 之后的 内容长度
	bodyLength := bl + decoder.LengthAdjustment

	if LogLevel <= 0 {
		slog.Info("读取到bodyLength", "bodyLength", bodyLength)
	}

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

	if LogLevel <= 0 {
		slog.Info("读取到BodyBytes", "BodyBytes", fmt.Sprintf("%x", result.BodyBytes))
	}

	result.FrameLength = len(result.HeaderBytes) + decoder.LengthFieldLength + len(result.BodyBytes)

	return result, nil
}
