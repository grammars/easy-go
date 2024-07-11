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
	BytesRead  chan int
	BytesWrite chan int
	ValidNum   chan int
	InvalidNum chan int
}

func (m *Monitor) Start() {
	m.BytesRead = make(chan int)
	m.BytesWrite = make(chan int)
	m.ValidNum = make(chan int)
	m.InvalidNum = make(chan int)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	nTick := 0
	sumBytesRead := 0
	sumBytesWrite := 0
	sumReadTimes := 0
	sumWriteTimes := 0
	sumValid := 0
	sumInvalid := 0
	for {
		select {
		case num := <-m.BytesRead:
			sumBytesRead += num
			sumReadTimes++
			break
		case num := <-m.BytesWrite:
			sumBytesWrite += num
			sumWriteTimes++
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
			fmt.Printf("%s 嘀嗒 %d 收到字节%d 次数%d 写出字节%d 次数%d 有效数=%d 失效数=%d \n",
				curTime.Format("2006-01-02 15:04:05"),
				nTick, sumBytesRead, sumReadTimes, sumBytesWrite, sumWriteTimes, sumValid, sumInvalid)
			break
		}
	}
}
