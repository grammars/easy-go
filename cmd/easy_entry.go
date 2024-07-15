package main

import (
	"flag"
	"fmt"
	"github.com/grammars/easy-go/file"
	"github.com/grammars/easy-go/socket"
	"path/filepath"
	"time"
	"unsafe"
)

func main() {
	runType := flag.String("run", "", "which sub command to run")
	host := flag.String("host", "192.168.11.11", "host to connect")
	port := flag.Int("port", 8181, "port to listen or connect")
	tls := flag.Bool("tls", false, "enable TLS")
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
	case "sws":
		RunSocketWebServer(*port, *tls)
		break
	case "swc":
		RunSocketWebClient(*host, *port, *tls, *nc)
		break
	default:
		RunDefault(*runType, *host, *port, *tls, *nc)
	}
}

func RunSocketRawServer(port int) {
	fmt.Printf("RunSocketRawServer port=%d\n", port)
	srv := &socket.RawServer{Port: port, Monitor: socket.CreateMonitorStart()}
	srv.Start()
}

func RunSocketRawClient(host string, port int, clientNum int) {
	fmt.Printf("RunSocketRawClient host=%s port=%d clientNum=%d\n", host, port, clientNum)
	socket.TestManyRawClient(host, port, clientNum)
	time.Sleep(30 * time.Minute)
}

func RunSocketWebServer(port int, tls bool) {
	fmt.Printf("RunSocketWebServer port=%d tls=%v\n", port, tls)
	srv := &socket.WebServer{Port: port, Monitor: socket.CreateMonitorStart(), TLS: tls}
	if tls {
		srv.CrtFile = filepath.Join(file.GetExeDir(), "server.crt")
		srv.KeyFile = filepath.Join(file.GetExeDir(), "server.key")
		fmt.Printf("CrtFile=%s KeyFile=%s\n", srv.CrtFile, srv.KeyFile)
	}
	srv.StartDefault()
}

func RunSocketWebClient(host string, port int, tls bool, clientNum int) {
	fmt.Printf("RunSocketWebClient host=%s port=%d tls=%v clientNum=%d\n", host, port, tls, clientNum)
	socket.TestManyWebClient(host, port, tls, clientNum)
	time.Sleep(30 * time.Minute)
}

func RunDefault(runType string, host string, port int, tls bool, nc int) {
	fmt.Printf("Unknown sub command:%s %s %d %v %d\n", runType, host, port, tls, nc)
	var arr []int
	fmt.Println("arr0 ", unsafe.Pointer(&arr))
	arr = append(arr, 1)
	fmt.Println("arr1 ", unsafe.Pointer(&arr))
	arr = append(arr, 2)
	fmt.Println("arr2 ", unsafe.Pointer(&arr))
	fmt.Println(arr)
	absPath, _ := filepath.Abs("文件.txt")
	fmt.Printf("absPath=%s\n", absPath)
	absExePath := file.ConvAbsPathRelExe("相对于执行文件的文件.doc")
	fmt.Printf("absExePath=%s\n", absExePath)
}
