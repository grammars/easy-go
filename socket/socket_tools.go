package socket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/grammars/easy-go/sugar"
	"log/slog"
	"net"
	"strings"
	"time"
)

func GetWsPath(path string) string {
	path = sugar.EnsureNotBlank(path, "/ws")
	if !strings.HasPrefix(path, "/") {
		return "/" + path
	}
	return path
}

var connCloseSum = 0

func CloseConn(conn net.Conn) {
	err := conn.Close()
	if err != nil {
		fmt.Printf("连接关闭失败，错误:%v\n\n", err)
	} else {
		connCloseSum++
		fmt.Printf("成功关闭1个连接 (累计关闭%d)\n", connCloseSum)
	}
}

func CloseWebConn(conn *websocket.Conn) {
	err := conn.Close()
	if err != nil {
		slog.Error("关闭webSocket连接 失败", err)
	} else {
		slog.Info("关闭webSocket连接 成功")
	}
}

type Monitor struct {
	BytesRead     chan int
	BytesWrite    chan int
	AcceptFailNum chan int
	ValidNum      chan int
	InvalidNum    chan int
	IntervalMilli int
}

func (m *Monitor) Start() {
	m.BytesRead = make(chan int)
	m.BytesWrite = make(chan int)
	m.AcceptFailNum = make(chan int)
	m.ValidNum = make(chan int)
	m.InvalidNum = make(chan int)
	if m.IntervalMilli == 0 {
		m.IntervalMilli = 3000
	}
	ticker := time.NewTicker(time.Millisecond * time.Duration(m.IntervalMilli))
	defer ticker.Stop()
	nTick := 0
	var snBytesRead = &SumNum[int64]{}
	var snBytesWrite = &SumNum[int64]{}
	sumAcceptFail := 0
	sumValid := 0
	sumInvalid := 0
	var lastTickTime = time.Now().UnixMilli()
	for {
		select {
		case num := <-m.BytesRead:
			snBytesRead.Add(int64(num))
			break
		case num := <-m.BytesWrite:
			snBytesWrite.Add(int64(num))
			break
		case num := <-m.AcceptFailNum:
			sumAcceptFail += num
			break
		case <-m.ValidNum:
			sumValid++
			break
		case <-m.InvalidNum:
			sumInvalid++
			break
		case <-ticker.C:
			nTick++
			curTime := time.Now()
			deltaTime := curTime.UnixMilli() - lastTickTime
			lastTickTime = time.Now().UnixMilli()
			fmt.Printf("%s 嘀嗒 %d 读取字节%d 次数%d 速度%s %s 写出字节%d 次数%d 速度%s %s 有效数=%d 失效数=%d accept失败数=%d \n",
				curTime.Format(time.DateTime), nTick,
				snBytesRead.Total, snBytesRead.Times, speed(snBytesRead.DeltaNum(), deltaTime, "B"), speed(int64(snBytesRead.DeltaTimes()), deltaTime, "条"),
				snBytesWrite.Total, snBytesWrite.Times, speed(snBytesWrite.DeltaNum(), deltaTime, "B"), speed(int64(snBytesWrite.DeltaTimes()), deltaTime, "条"),
				sumValid, sumInvalid, sumAcceptFail)
			break
		}
	}
}

func speed(bs int64, deltaTime int64, unit string) string {
	if deltaTime <= 0 {
		return "--"
	}
	bs = bs * 1000 / deltaTime
	if bs < 1024 {
		return fmt.Sprintf("%d%s", bs, unit)
	} else if bs < 1024*1024 {
		return fmt.Sprintf("%.2fK%s", float64(bs)/1024, unit)
	} else if bs < 1024*1024*1024 {
		return fmt.Sprintf("%.2fM%s", float64(bs)/1024/1024, unit)
	}
	return fmt.Sprintf("%.2fG%s", float64(bs)/1024/1024/1024, unit)
}

type SumNum[T int | int64] struct {
	Total     T
	Times     int
	LastNum   T
	LastTimes int
}

func (sn *SumNum[T]) Add(num T) {
	sn.Total += num
	sn.Times++
}

func (sn *SumNum[T]) DeltaNum() T {
	d := sn.Total - sn.LastNum
	sn.LastNum = sn.Total
	return d
}

func (sn *SumNum[T]) DeltaTimes() int {
	d := sn.Times - sn.LastTimes
	sn.LastTimes = sn.Times
	return d
}
