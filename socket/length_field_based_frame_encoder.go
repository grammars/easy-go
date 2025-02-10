package socket

import (
	"encoding/binary"
	"fmt"
	"github.com/grammars/easy-go/tool"
)

type LengthFieldType int

const (
	lengthFieldTypeUint8  LengthFieldType = 1
	lengthFieldTypeUint16 LengthFieldType = 2 // 1
	lengthFieldTypeUint32 LengthFieldType = 4 // 2
)

type LengthFieldBasedFrameEncoder struct {
	ByteOrder         binary.ByteOrder
	MaxFrameLength    int
	LengthFieldOffset int
	LengthFieldLength LengthFieldType
}

func (encoder *LengthFieldBasedFrameEncoder) Encode(headerBytes, bodyBytes []byte) ([]byte, error) {
	ba := tool.NewByteArray()
	if encoder.LengthFieldOffset > 0 {
		if headerBytes != nil && len(headerBytes) > 0 {
			err := ba.WriteBytes(headerBytes[0:min(encoder.LengthFieldOffset, len(headerBytes))])
			if err != nil {
				return nil, err
			}
		}
	}
	var bodyLength = 0
	if bodyBytes != nil {
		bodyLength = len(bodyBytes)
	}
	blb := make([]byte, encoder.LengthFieldLength)
	if encoder.LengthFieldLength == lengthFieldTypeUint8 {
		blb[0] = uint8(bodyLength)
	} else if encoder.LengthFieldLength == lengthFieldTypeUint16 {
		encoder.ByteOrder.PutUint16(blb, uint16(bodyLength))
	} else if encoder.LengthFieldLength == lengthFieldTypeUint32 {
		encoder.ByteOrder.PutUint32(blb, uint32(bodyLength))
	} else {
		return nil, fmt.Errorf("LengthFieldBasedFrameEncoder encode error: LengthFieldLength is invalid")
	}
	err := ba.WriteBytes(blb)
	if err != nil {
		return nil, err
	}
	if bodyLength > 0 {
		err = ba.WriteBytes(bodyBytes)
		if err != nil {
			return nil, err
		}
	}
	return ba.Bytes(), nil
}
