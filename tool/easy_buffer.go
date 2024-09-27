package tool

import (
	"bytes"
	"encoding/binary"
)

type ByteArray struct {
	buf       *bytes.Buffer
	byteOrder binary.ByteOrder
}

func NewByteArray() *ByteArray {
	return &ByteArray{buf: &bytes.Buffer{}, byteOrder: binary.BigEndian}
}

func (ins *ByteArray) Order(byteOrder binary.ByteOrder) *ByteArray {
	ins.byteOrder = byteOrder
	return ins
}

func (ins *ByteArray) WriteByte(value byte) error {
	return ins.buf.WriteByte(value)
}

func (ins *ByteArray) ReadByte() (byte, error) {
	return ins.buf.ReadByte()
}

// WriteUint8 type byte = uint8
func (ins *ByteArray) WriteUint8(value uint8) error {
	return ins.buf.WriteByte(value)
}

func (ins *ByteArray) ReadUint8() (uint8, error) {
	b, err := ins.buf.ReadByte()
	return b, err
}

func (ins *ByteArray) WriteInt8(value int8) error {
	return ins.buf.WriteByte(byte(value))
}

func (ins *ByteArray) ReadInt8() (int8, error) {
	b, err := ins.buf.ReadByte()
	return int8(b), err
}

func __ReadTemplate__[T any](ins *ByteArray) (T, error) {
	var value T
	err := binary.Read(ins.buf, ins.byteOrder, &value)
	if err != nil {
		return value, err
	}
	return value, nil
}

func (ins *ByteArray) WriteUint16(value uint16) error {
	bs := make([]byte, 2)
	ins.byteOrder.PutUint16(bs, value)
	_, err := ins.buf.Write(bs)
	return err
}

func (ins *ByteArray) ReadUint16() (uint16, error) {
	return __ReadTemplate__[uint16](ins)
}

func (ins *ByteArray) WriteInt16(value int16) error {
	bs := make([]byte, 2)
	ins.byteOrder.PutUint16(bs, uint16(value))
	_, err := ins.buf.Write(bs)
	return err
}

func (ins *ByteArray) ReadInt16() (int16, error) {
	return __ReadTemplate__[int16](ins)
}

func (ins *ByteArray) WriteUint32(value uint32) error {
	bs := make([]byte, 4)
	ins.byteOrder.PutUint32(bs, value)
	_, err := ins.buf.Write(bs)
	return err
}

func (ins *ByteArray) ReadUint32() (uint32, error) {
	return __ReadTemplate__[uint32](ins)
}

func (ins *ByteArray) WriteInt32(value int32) error {
	bs := make([]byte, 4)
	ins.byteOrder.PutUint32(bs, uint32(value))
	_, err := ins.buf.Write(bs)
	return err
}

func (ins *ByteArray) ReadInt32() (int32, error) {
	return __ReadTemplate__[int32](ins)
}

func (ins *ByteArray) WriteUint64(value uint64) error {
	bs := make([]byte, 8)
	ins.byteOrder.PutUint64(bs, value)
	_, err := ins.buf.Write(bs)
	return err
}

func (ins *ByteArray) ReadUint64() (uint64, error) {
	return __ReadTemplate__[uint64](ins)
}

func (ins *ByteArray) WriteInt64(value int64) error {
	bs := make([]byte, 8)
	ins.byteOrder.PutUint64(bs, uint64(value))
	_, err := ins.buf.Write(bs)
	return err
}

func (ins *ByteArray) ReadInt64() (int64, error) {
	return __ReadTemplate__[int64](ins)
}

func (ins *ByteArray) WriteFloat32(value float32) error {
	return binary.Write(ins.buf, ins.byteOrder, value)
}

func (ins *ByteArray) ReadFloat32() (float32, error) {
	return __ReadTemplate__[float32](ins)
}

func (ins *ByteArray) WriteFloat64(value float64) error {
	return binary.Write(ins.buf, ins.byteOrder, value)
}

func (ins *ByteArray) ReadFloat64() (float64, error) {
	return __ReadTemplate__[float64](ins)
}

func (ins *ByteArray) WriteBool(value bool) error {
	if value {
		return ins.buf.WriteByte(1)
	}
	return ins.buf.WriteByte(0)
}

func (ins *ByteArray) ReadBool() (bool, error) {
	b, err := ins.buf.ReadByte()
	if err != nil {
		return false, err
	}
	if b == 0 {
		return false, nil
	}
	return true, nil
}
