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
	InvalidNum chan int
}

func (m *Monitor) Start() {
	m.BytesRead = make(chan int)
	m.InvalidNum = make(chan int)
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	nTick := 0
	sumBytesRead := 0
	sumReadTimes := 0
	sumInvalid := 0
	for {
		select {
		case num := <-m.BytesRead:
			sumBytesRead += num
			sumReadTimes++
			break
		case <-m.InvalidNum:
			sumInvalid++
			break
		case <-ticker.C:
			nTick++
			fmt.Printf("嘀嗒 %d 收到字节%d 次数%d 失效数=%d \n", nTick, sumBytesRead, sumReadTimes, sumInvalid)
			break
		}
	}
}
