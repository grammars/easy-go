package main

import (
	"flag"
	"fmt"
	"github.com/grammars/easy-go/socket"
	"time"
	"unsafe"
)

func main() {
	runType := flag.String("run", "", "which sub command to run")
	host := flag.String("host", "192.168.11.11", "host to connect")
	port := flag.Int("port", 8181, "port to listen or connect")
	nc := flag.Int("nc", 1, "number of client")
	//fmt.Printf("即将解析 runType=%s nc=%d\n", *runType, *nc)
	flag.Parse()
	//fmt.Printf("解析之后 runType=%s nc=%d\n", *runType, *nc)
	switch *runType {
	case "srs":
		RunSocketRawServer(*port)
		break
	case "src":
		RunSocketRawClient(*host, *port, *nc)
		break
	default:
		RunDefault(*runType)
	}
}

func RunSocketRawServer(port int) {
	fmt.Printf("RunSocketRawServer port=%d\n", port)
	socket.CreateRawServer(port).Start()
}

func RunSocketRawClient(host string, port int, clientNum int) {
	fmt.Printf("RunSocketRawClient host=%s port=%d clientNum=%d\n", host, port, clientNum)
	socket.TestManyRawClient(host, port, clientNum)
	time.Sleep(30 * time.Minute)
}

func RunDefault(runType string) {
	fmt.Println("Unknown sub command:", runType)
	var arr []int
	fmt.Println("arr0 ", unsafe.Pointer(&arr))
	arr = append(arr, 1)
	fmt.Println("arr1 ", unsafe.Pointer(&arr))
	arr = append(arr, 2)
	fmt.Println("arr2 ", unsafe.Pointer(&arr))
	fmt.Println(arr)
}
