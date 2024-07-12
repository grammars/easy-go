package socket

import (
	"fmt"
	"net"
	"time"
)

func CloseConn(conn net.Conn) {
	err := conn.Close()
	if err != nil {
		fmt.Printf("连接关闭失败，错误:%v\n\n", err)
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
	var sumBytesRead int64 = 0
	var sumBytesWrite int64 = 0
	sumReadTimes := 0
	sumWriteTimes := 0
	sumAcceptFail := 0
	sumValid := 0
	sumInvalid := 0
	var lastTickTime = time.Now().UnixMilli()
	var lastBytesRead int64 = 0
	var lastBytesWrite int64 = 0
	for {
		select {
		case num := <-m.BytesRead:
			sumBytesRead += int64(num)
			sumReadTimes++
			break
		case num := <-m.BytesWrite:
			sumBytesWrite += int64(num)
			sumWriteTimes++
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
			deltaBytesRead := sumBytesRead - lastBytesRead
			deltaBytesWrite := sumBytesWrite - lastBytesWrite
			fmt.Printf("%s 嘀嗒 %d 收到字节%d 次数%d 速度%s 写出字节%d 次数%d 速度%s 有效数=%d 失效数=%d accept失败数=%d \n",
				curTime.Format(time.DateTime), nTick,
				sumBytesRead, sumReadTimes, speed(deltaBytesRead, deltaTime),
				sumBytesWrite, sumWriteTimes, speed(deltaBytesWrite, deltaTime),
				sumValid, sumInvalid, sumAcceptFail)
			lastBytesRead = sumBytesRead
			lastBytesWrite = sumBytesWrite
			break
		}
	}
}

func speed(bs, deltaTime int64) string {
	if deltaTime <= 0 {
		return "--"
	}
	bs = bs * 1000 / deltaTime
	if bs < 1024 {
		return fmt.Sprintf("%dB", bs)
	} else if bs < 1024*1024 {
		return fmt.Sprintf("%.2fKB", float64(bs)/1024)
	} else if bs < 1024*1024*1024 {
		return fmt.Sprintf("%.2fMB", float64(bs)/1024/1024)
	}
	return fmt.Sprintf("%.2fGB", float64(bs)/1024/1024/1024)
}
