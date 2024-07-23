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
}

// VisitorServer Visitor所属的server
type VisitorServer interface {
	GetStartTime() time.Time
}

type Visitor[VD any] struct {
	index uint64
	uid   uint64
	Conn  VisitorConnection
	Data  *VD
}

type VisitorServerHandler[VD any] interface {
	OnConnect(visitor *Visitor[VD])
	OnMessage(visitor *Visitor[VD], message *[]byte)
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
	visitor.uid = uint64(vm.server.GetStartTime().UnixMilli()) + visitor.index
	vm.holder.Store(visitor.uid, visitor)
	return visitor
}

func (vm *VisitorMap[VD]) Remove(visitorUid uint64) *Visitor[VD] {
	visitor, ok := vm.holder.Load(visitorUid)
	if ok {
		vm.holder.Delete(visitorUid)
	}

	return visitor.(*Visitor[VD])
}

func (vm *VisitorMap[VD]) Print() {
	vm.holder.Range(func(k, v any) bool {
		uid := k.(uint64)
		visitor := v.(*Visitor[VD])
		slog.Info("打印visitorMap", "uid", uid, "index", visitor.index, "addr", visitor.Conn.RemoteAddr())
		return true
	})
}

func CreateVisitorMap[VD any](server VisitorServer) *VisitorMap[VD] {
	vm := &VisitorMap[VD]{server: server, history: 0}
	vm.holder = new(sync.Map)
	return vm
}
