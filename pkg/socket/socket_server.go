package msock

import (
	"bufio"
	"fmt"
	"io"
	"net"
)

func ReadWriteAsServer(conn net.Conn, monitor *Monitor) {
	defer CloseConn(conn)
	reader := bufio.NewReader(conn)
	for {
		var buf [1024]byte
		n, err := reader.Read(buf[:])
		if err != nil && err != io.EOF {
			//fmt.Printf("读取失败,err:%v\n", err)
			monitor.InvalidNum <- 1
			break
		}
		monitor.BytesRead <- n
		got := string(buf[:n])
		fmt.Println("接收到的数据: ", got)
		nWrite, err := conn.Write([]byte("收到了：" + got))
		if err != nil {
			return
		}
		fmt.Println("回复:nWrite=", nWrite)
	}
}
