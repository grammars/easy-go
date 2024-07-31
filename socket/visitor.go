package socket

import (
	"log/slog"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type VisitorConnection interface {
	RemoteAddr() net.Addr
	Write(b []byte) (n int, err error)
	// WriteSafe 协程安全的写入
	WriteSafe(b []byte) (n int, err error)
}

// VisitorServer Visitor所属的server
type VisitorServer interface {
	GetStartTime() time.Time
}

type Visitor[VD any] struct {
	index uint64
	Uid   uint64
	Conn  VisitorConnection
	Data  *VD
	Ext   any
}

type VisitorServerHandler[VD any] interface {
	OnConnect(visitor *Visitor[VD])

	// OnMessage 当接收到消息时
	//
	// messageType 消息类型, 当用于WebSocket的时候, 采用 https://www.rfc-editor.org/rfc/rfc6455.html#section-11.8 标准中的枚举;
	// 自定义协议时, 暂不约定类型
	OnMessage(visitor *Visitor[VD], message CodecResult, messageType int)
	OnDisconnect(visitor *Visitor[VD])
}

type VisitorMap[VD any] struct {
	history uint64
	server  VisitorServer
	holder  *sync.Map //map[uint64]*Visitor
}

func (vm *VisitorMap[VD]) Append(conn VisitorConnection) *Visitor[VD] {
	visitor := &Visitor[VD]{Conn: conn}
	visitor.index = atomic.AddUint64(&vm.history, 1)
	visitor.Uid = uint64(vm.server.GetStartTime().UnixMilli()) + visitor.index
	vm.holder.Store(visitor.Uid, visitor)
	return visitor
}

func (vm *VisitorMap[VD]) Remove(visitorUid uint64) *Visitor[VD] {
	visitor, ok := vm.holder.Load(visitorUid)
	if ok {
		vm.holder.Delete(visitorUid)
	}

	return visitor.(*Visitor[VD])
}

func (vm *VisitorMap[VD]) Get(visitorUid uint64) *Visitor[VD] {
	visitor, ok := vm.holder.Load(visitorUid)
	if ok {
		return visitor.(*Visitor[VD])
	}
	return nil
}

func (vm *VisitorMap[VD]) Print() {
	vm.holder.Range(func(k, v any) bool {
		uid := k.(uint64)
		visitor := v.(*Visitor[VD])
		slog.Info("打印visitorMap", "Uid", uid, "index", visitor.index, "addr", visitor.Conn.RemoteAddr())
		return true
	})
}

func CreateVisitorMap[VD any](server VisitorServer) *VisitorMap[VD] {
	vm := &VisitorMap[VD]{server: server, history: 0}
	vm.holder = new(sync.Map)
	return vm
}
