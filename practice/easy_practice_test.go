package practice

import (
	"fmt"
	"testing"
	"time"
	"unsafe"
)

func TestSth(t *testing.T) {
	fmt.Println("TestSth")
	var app0 = Apple{Name: "黑富士", weight: 4.36}
	addr0 := unsafe.Pointer(&app0)
	fmt.Printf("app0=%v addr0=%p\n", app0, addr0)

	app1 := CreateApple("红富士", 3.666)
	addr1 := unsafe.Pointer(app1)
	fmt.Printf("app1=%v addr1=%p\n", app1, addr1)
	app1.Name = "新富士"
	time.Sleep(3 * time.Second)
}

type Apple struct {
	Name   string
	weight float32
}

func CreateApple(name string, weight float32) *Apple {
	app := Apple{Name: name, weight: weight}
	addr := unsafe.Pointer(&app)
	fmt.Printf("创建方法中的 app=%v addr=%p\n", app, addr)
	go func() {
		fmt.Println("进入协程")
		addrInRoute := unsafe.Pointer(&app)
		fmt.Printf("协程方法中的 app=%v addr=%p\n", app, addrInRoute)
		time.Sleep(2 * time.Second)
		fmt.Printf("2秒后 app=%v \n", app)
	}()
	return &app
}
