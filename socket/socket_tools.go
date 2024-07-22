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
	Name          string
	IntervalMilli int
	BytesRead     chan int
	BytesWrite    chan int
	AcceptFailNum chan int
	ValidNum      chan int
	InvalidNum    chan int
	Stat          *MonitorStat
}

type MonitorStat struct {
	snBytesRead   *SumNum[int64]
	snBytesWrite  *SumNum[int64]
	sumAcceptFail int
	sumValid      int
	sumInvalid    int
}

func (stat *MonitorStat) Reset() {
	stat.snBytesRead = &SumNum[int64]{}
	stat.snBytesWrite = &SumNum[int64]{}
	stat.sumAcceptFail = 0
	stat.sumValid = 0
	stat.sumInvalid = 0
}

func (stat *MonitorStat) ToMap() map[string]any {
	m := make(map[string]any)
	m["bytesReadTotal"] = stat.snBytesRead.Total
	m["bytesWriteTotal"] = stat.snBytesWrite.Total
	m["sumValid"] = stat.sumValid
	m["sumInvalid"] = stat.sumInvalid
	m["online"] = stat.sumValid - stat.sumInvalid
	return m
}

func CreateMonitorStart(name string, intervalMilli int) *Monitor {
	monitor := &Monitor{Name: name, IntervalMilli: intervalMilli}
	go monitor.Start()
	return monitor
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
	m.Stat = &MonitorStat{}
	m.Stat.Reset()
	ticker := time.NewTicker(time.Millisecond * time.Duration(m.IntervalMilli))
	defer ticker.Stop()
	nTick := 0
	var lastTickTime = time.Now().UnixMilli()
	for {
		select {
		case num := <-m.BytesRead:
			m.Stat.snBytesRead.Add(int64(num))
			break
		case num := <-m.BytesWrite:
			m.Stat.snBytesWrite.Add(int64(num))
			break
		case num := <-m.AcceptFailNum:
			m.Stat.sumAcceptFail += num
			break
		case <-m.ValidNum:
			m.Stat.sumValid++
			break
		case <-m.InvalidNum:
			m.Stat.sumInvalid++
			break
		case <-ticker.C:
			nTick++
			curTime := time.Now()
			deltaTime := curTime.UnixMilli() - lastTickTime
			lastTickTime = time.Now().UnixMilli()
			onlineNum := m.Stat.sumValid - m.Stat.sumInvalid
			fmt.Printf("[%s] %s 嘀嗒 %d 读取字节%d 次数%d 速度%s %s 写出字节%d 次数%d 速度%s %s 在线=%d 有效数=%d 失效数=%d accept失败数=%d \n",
				sugar.EnsureNotBlank(m.Name, "默认监视器"),
				curTime.Format(time.DateTime), nTick,
				m.Stat.snBytesRead.Total, m.Stat.snBytesRead.Times, speed(m.Stat.snBytesRead.DeltaNum(), deltaTime, "B"), speed(int64(m.Stat.snBytesRead.DeltaTimes()), deltaTime, "条"),
				m.Stat.snBytesWrite.Total, m.Stat.snBytesWrite.Times, speed(m.Stat.snBytesWrite.DeltaNum(), deltaTime, "B"), speed(int64(m.Stat.snBytesWrite.DeltaTimes()), deltaTime, "条"),
				onlineNum, m.Stat.sumValid, m.Stat.sumInvalid, m.Stat.sumAcceptFail)
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
