package socket

import (
	"fmt"
	"net"
	"time"
)

type Utils struct{}

func (u *Utils) CreateRawServer(port int) *RawServer {
	return &RawServer{Port: port}
}

func (u *Utils) CreateRawClient(addr string, port int, name string) *RawClient {
	return &RawClient{Addr: addr, Port: port, Name: name}
}

type RawServer struct {
	Port int
}

func (srv *RawServer) Start() {
	fmt.Println("开始启动SocketServer")

	monitor := Monitor{}
	go monitor.Start()

	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", srv.Port))
	if err != nil {
		fmt.Printf("启动监听失败，错误:%v\n\n", err)
		panic(err)
	}
	var conn net.Conn
	defer CloseConn(conn)
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Printf("接收连接失败:%v\n", err)
			continue
		}
		fmt.Printf("有一个客户端连接我成功了，来自:%v\n", conn.RemoteAddr())
		go ReadWriteAsServer(conn, &monitor)
	}
}

type RawClient struct {
	Addr string
	Port int
	Name string
}

func (cli *RawClient) Start() {
	fmt.Printf("RawClient %s Start!\n", cli.Name)
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", cli.Addr, cli.Port))
	if err != nil {
		fmt.Printf("[%s]创建连接失败，错误:%v\n", cli.Name, err)
		return
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("连接关闭失败，错误:%v\n\n", err)
		}
	}(conn)
	fmt.Printf("[%s]连接服务端成功:%v\n", cli.Name, conn.RemoteAddr())

	for i := 0; i < 100; i++ {
		_, err := conn.Write([]byte(fmt.Sprintf("[%s]hello world(%d)", cli.Name, i)))
		if err != nil {
			fmt.Printf("发送消息失败:%v\n", err)
			return
		}
		time.Sleep(5000 * time.Millisecond)
		var buf [1024]byte
		n, err := conn.Read(buf[:])
		if err != nil {
			fmt.Printf("read failed, err:%v\n", err)
			return
		}
		if n > 0 {
			fmt.Printf("[%s] 收到服务端回复:%s\n", cli.Name, string(buf[:n]))
		} else {
			fmt.Printf("[%s] 没有收到服务端回复:\n", cli.Name)
		}
	}
}
