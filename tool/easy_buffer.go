package tool

import (
	"bytes"
	"encoding/binary"
)

type ByteArray struct {
	buf        *bytes.Buffer
	byteOrder  binary.ByteOrder
	length     int // 内容总长度
	readLength int // 已读长度
}

func NewByteArray() *ByteArray {
	return &ByteArray{buf: &bytes.Buffer{}, byteOrder: binary.BigEndian}
}

func (ins *ByteArray) Order(byteOrder binary.ByteOrder) *ByteArray {
	ins.byteOrder = byteOrder
	return ins
}

func (ins *ByteArray) Buffer() *bytes.Buffer {
	return ins.buf
}

func (ins *ByteArray) Bytes() []byte {
	return ins.buf.Bytes()
}

// Length 字节内容总长度
func (ins *ByteArray) Length() int {
	return ins.length
}

// Available 剩余可读长度
func (ins *ByteArray) Available() int {
	return ins.length - ins.readLength
}

func (ins *ByteArray) AfterRead(err error, n int) error {
	if err == nil {
		ins.length += n
	}
	return err
}

func (ins *ByteArray) WriteByte(value byte) error {
	return ins.AfterRead(ins.buf.WriteByte(value), 1)
}

func (ins *ByteArray) ReadByte() (byte, error) {
	v, err := ins.buf.ReadByte()
	if err == nil {
		ins.readLength += 1
	}
	return v, err
}

// WriteUint8 type byte = uint8
func (ins *ByteArray) WriteUint8(value uint8) error {
	return ins.AfterRead(ins.buf.WriteByte(value), 1)
}

func (ins *ByteArray) ReadUint8() (uint8, error) {
	b, err := ins.buf.ReadByte()
	if err == nil {
		ins.readLength += 1
	}
	return b, err
}

func (ins *ByteArray) WriteInt8(value int8) error {
	return ins.AfterRead(ins.buf.WriteByte(byte(value)), 1)
}

func (ins *ByteArray) ReadInt8() (int8, error) {
	b, err := ins.buf.ReadByte()
	if err == nil {
		ins.readLength += 1
	}
	return int8(b), err
}

func __ReadTemplate__[T any](ins *ByteArray, nRead int) (T, error) {
	var value T
	err := binary.Read(ins.buf, ins.byteOrder, &value)
	if err != nil {
		return value, err
	}
	ins.readLength += nRead
	return value, nil
}

func (ins *ByteArray) WriteUint16(value uint16) error {
	bs := make([]byte, 2)
	ins.byteOrder.PutUint16(bs, value)
	_, err := ins.buf.Write(bs)
	return ins.AfterRead(err, 2)
}

func (ins *ByteArray) ReadUint16() (uint16, error) {
	return __ReadTemplate__[uint16](ins, 2)
}

func (ins *ByteArray) WriteInt16(value int16) error {
	bs := make([]byte, 2)
	ins.byteOrder.PutUint16(bs, uint16(value))
	_, err := ins.buf.Write(bs)
	return ins.AfterRead(err, 2)
}

func (ins *ByteArray) ReadInt16() (int16, error) {
	return __ReadTemplate__[int16](ins, 2)
}

func (ins *ByteArray) WriteUint32(value uint32) error {
	bs := make([]byte, 4)
	ins.byteOrder.PutUint32(bs, value)
	_, err := ins.buf.Write(bs)
	return ins.AfterRead(err, 4)
}

func (ins *ByteArray) ReadUint32() (uint32, error) {
	return __ReadTemplate__[uint32](ins, 4)
}

func (ins *ByteArray) WriteInt32(value int32) error {
	bs := make([]byte, 4)
	ins.byteOrder.PutUint32(bs, uint32(value))
	_, err := ins.buf.Write(bs)
	return ins.AfterRead(err, 4)
}

func (ins *ByteArray) ReadInt32() (int32, error) {
	return __ReadTemplate__[int32](ins, 4)
}

func (ins *ByteArray) WriteUint64(value uint64) error {
	bs := make([]byte, 8)
	ins.byteOrder.PutUint64(bs, value)
	_, err := ins.buf.Write(bs)
	return ins.AfterRead(err, 8)
}

func (ins *ByteArray) ReadUint64() (uint64, error) {
	return __ReadTemplate__[uint64](ins, 8)
}

func (ins *ByteArray) WriteInt64(value int64) error {
	bs := make([]byte, 8)
	ins.byteOrder.PutUint64(bs, uint64(value))
	_, err := ins.buf.Write(bs)
	return ins.AfterRead(err, 8)
}

func (ins *ByteArray) ReadInt64() (int64, error) {
	return __ReadTemplate__[int64](ins, 8)
}

func (ins *ByteArray) WriteFloat32(value float32) error {
	return ins.AfterRead(binary.Write(ins.buf, ins.byteOrder, value), 4)
}

func (ins *ByteArray) ReadFloat32() (float32, error) {
	return __ReadTemplate__[float32](ins, 4)
}

func (ins *ByteArray) WriteFloat64(value float64) error {
	return ins.AfterRead(binary.Write(ins.buf, ins.byteOrder, value), 8)
}

func (ins *ByteArray) ReadFloat64() (float64, error) {
	return __ReadTemplate__[float64](ins, 8)
}

func (ins *ByteArray) WriteBool(value bool) error {
	var err error
	if value {
		err = ins.buf.WriteByte(1)
	} else {
		err = ins.buf.WriteByte(0)
	}
	return ins.AfterRead(err, 1)
}

func (ins *ByteArray) ReadBool() (bool, error) {
	b, err := ins.buf.ReadByte()
	if err != nil {
		return false, err
	}
	ins.readLength += 1
	if b == 0 {
		return false, nil
	}
	return true, nil
}

func (ins *ByteArray) WriteInt(value int) error {
	bs := make([]byte, 4)
	ins.byteOrder.PutUint32(bs, uint32(value))
	_, err := ins.buf.Write(bs)
	return ins.AfterRead(err, 4)
}

func (ins *ByteArray) ReadInt() (int, error) {
	u, err := __ReadTemplate__[uint32](ins, 4)
	return int(u), err
}
