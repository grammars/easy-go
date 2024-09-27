package tool

import "testing"

func TestByteArray(t *testing.T) {
	ba := NewByteArray()
	// byte
	err := ba.WriteByte(235)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("累计写入=%d字节 剩余可读=%d字节", ba.Length(), ba.Available())
	vByte, _ := ba.ReadByte()
	t.Logf("vByte=%d", vByte)

	// uint8
	err = ba.WriteUint8(255)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("累计写入=%d字节 剩余可读=%d字节", ba.Length(), ba.Available())
	vUint8, _ := ba.ReadUint8()
	t.Logf("vUint8=%d", vUint8)

	// int8
	err = ba.WriteInt8(-123)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("累计写入=%d字节 剩余可读=%d字节", ba.Length(), ba.Available())
	vInt8, _ := ba.ReadInt8()
	t.Logf("vInt8=%d", vInt8)

	// uint16
	err = ba.WriteUint16(12345)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("累计写入=%d字节 剩余可读=%d字节", ba.Length(), ba.Available())
	vUint16, _ := ba.ReadUint16()
	t.Logf("vUint16=%d", vUint16)

	// int16
	err = ba.WriteInt16(-10010)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("累计写入=%d字节 剩余可读=%d字节", ba.Length(), ba.Available())
	vInt16, _ := ba.ReadInt16()
	t.Logf("vInt16=%d", vInt16)

	// uint32
	err = ba.WriteUint32(1234578910)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("累计写入=%d字节 剩余可读=%d字节", ba.Length(), ba.Available())
	vUint32, _ := ba.ReadUint32()
	t.Logf("vUint32=%d", vUint32)

	// int32
	err = ba.WriteInt32(-1008610010)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("累计写入=%d字节 剩余可读=%d字节", ba.Length(), ba.Available())
	vInt32, _ := ba.ReadInt32()
	t.Logf("vInt32=%d", vInt32)

	// uint64
	err = ba.WriteUint64(9876543210123456)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("累计写入=%d字节 剩余可读=%d字节", ba.Length(), ba.Available())
	vUint64, _ := ba.ReadUint64()
	t.Logf("vUint64=%d", vUint64)

	// int64
	err = ba.WriteInt64(-9876543210123456)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("累计写入=%d字节 剩余可读=%d字节", ba.Length(), ba.Available())
	vInt64, _ := ba.ReadInt64()
	t.Logf("vInt64=%d", vInt64)

	// float32
	err = ba.WriteFloat32(3.1415)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("累计写入=%d字节 剩余可读=%d字节", ba.Length(), ba.Available())
	vFloat32, _ := ba.ReadFloat32()
	t.Logf("float32=%f", vFloat32)

	// float64
	err = ba.WriteFloat64(-123456789.987654321)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("累计写入=%d字节 剩余可读=%d字节", ba.Length(), ba.Available())
	vFloat64, _ := ba.ReadFloat64()
	t.Logf("vFloat64=%f", vFloat64)

	// bool
	err = ba.WriteBool(true)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("累计写入=%d字节 剩余可读=%d字节", ba.Length(), ba.Available())
	vBool, _ := ba.ReadBool()
	t.Logf("vBool=%v", vBool)

	// int32
	err = ba.WriteInt(-989899880)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("累计写入=%d字节 剩余可读=%d字节", ba.Length(), ba.Available())
	vInt, _ := ba.ReadInt()
	t.Logf("vInt=%d", vInt)

	// stringUint8
	err = ba.WriteStringUint8("免费贴膜哦")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("累计写入=%d字节 剩余可读=%d字节", ba.Length(), ba.Available())
	vStr8, _ := ba.ReadStringUint8()
	t.Logf("vStr8=%s", vStr8)

	// stringUint16
	err = ba.WriteStringUint16("收费大宝剑")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("累计写入=%d字节 剩余可读=%d字节", ba.Length(), ba.Available())
	vStr16, _ := ba.ReadStringUint16()
	t.Logf("vStr16=%s", vStr16)

	// stringUint32
	err = ba.WriteStringUint32("山姆会员店")
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("累计写入=%d字节 剩余可读=%d字节", ba.Length(), ba.Available())
	vStr32, _ := ba.ReadStringUint32()
	t.Logf("vStr32=%s", vStr32)

	t.Logf("全部完成 累计写入=%d字节 剩余可读=%d字节", ba.Length(), ba.Available())
}
