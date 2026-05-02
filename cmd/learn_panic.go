package main

import (
	"fmt"
	"runtime"
	"time"
)

func LearnPanic() {
	fmt.Println("learn panic begin")
	SafeCall()
	time.Sleep(1 * time.Second)
	fmt.Println("learn panic end")
}

func SafeCall() {
	fmt.Println("safe call begin")
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("recover", r)
			// 打印完整的调用堆栈
			buf := make([]byte, 1024*1024) // 1MB 缓冲区
			n := runtime.Stack(buf, false)
			fmt.Printf("调用堆栈:\n%s\n", buf[:n])
		}
	}()
	doSomethingRisky("divide_by_zero", 0)
	fmt.Println("safe call end")
}

// doSomethingRisky 模拟一个可能触发 panic 的函数
func doSomethingRisky(action string, index int) {
	fmt.Printf("🚀 开始执行危险操作: %s\n", action)

	switch action {
	case "divide_by_zero":
		// 场景 1: 除以零 - 运行时错误
		fmt.Println("  -> 准备进行数学计算...")
		result := 10 / index // 如果 index 为 0，这里会触发 panic
		fmt.Printf("  -> 计算结果: %d\n", result)

	case "slice_out_of_bounds":
		// 场景 2: 切片越界 - 运行时错误
		fmt.Println("  -> 准备访问切片数据...")
		data := []string{"apple", "banana", "cherry"}
		// 如果 index 超出 [0, 2] 的范围，这里会触发 panic
		item := data[index]
		fmt.Printf("  -> 获取到的数据: %s\n", item)

	case "nil_pointer":
		// 场景 3: 空指针引用 - 程序员错误
		fmt.Println("  -> 准备操作一个空指针...")
		var p *int = nil
		// 解引用一个 nil 指针会立即触发 panic
		*p = 100

	case "manual_panic":
		// 场景 4: 手动触发 panic - 业务逻辑错误
		fmt.Println("  -> 检查业务逻辑...")
		if index < 0 {
			// 当遇到无法恢复的业务错误时，主动触发 panic
			panic(fmt.Sprintf("业务逻辑错误: 索引不能为负数 (%d)", index))
		}
		fmt.Printf("  -> 业务检查通过，索引为: %d\n", index)

	default:
		fmt.Println("  -> 未知操作，安全退出。")
	}

	fmt.Printf("✅ 危险操作 '%s' 成功完成！\n", action)
}
