package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello Gopher!")
	for i := 1; i <= 500; i++ {
		go PlayRawSocket(i)
	}
	time.Sleep(60 * time.Second)
}
