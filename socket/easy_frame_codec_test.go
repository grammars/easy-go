package socket

import (
	"encoding/binary"
	"testing"
)

func TestLengthFieldBasedFrame(t *testing.T) {
	t.Log("TestLengthFieldBasedFrame")
	order := binary.BigEndian
	lfl := lengthFieldTypeUint32
	encoder := &LengthFieldBasedFrameEncoder{ByteOrder: order, MaxFrameLength: 128, LengthFieldOffset: 2, LengthFieldLength: lfl}
	decoder := &LengthFieldBasedFrameDecoder{ByteOrder: order, MaxFrameLength: 128, LengthFieldOffset: 2, LengthFieldLength: int(lfl), LengthAdjustment: 0}

	messageBytes, err := encoder.Encode([]byte{11, 22, 33}, []byte{7, 8, 9})
	if err != nil {
		t.Error("Encode error:", err)
		return
	}
	t.Log("Encode message:", messageBytes)

	cr, err := decoder.ReadBytes(messageBytes)
	if err != nil {
		t.Error("ReadBytes error:", err)
		return
	}
	t.Log("Read bytes:", cr)
}
