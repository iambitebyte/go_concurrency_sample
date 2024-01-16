package main

import "fmt"

func main() {
	// 创建一个带缓冲的 channel，缓冲大小为 1
	ch := make(chan int, 1)

	// 启动一个 goroutine 发送数据到 channel
	go func() {
		ch <- 42
		fmt.Println("Sent data:", 42)
		ch <- 43
		fmt.Println("Sent data:", 43)
	}()

	go func() {
		ch <- 44
		fmt.Println("Sent data:", 44)
		ch <- 45
		fmt.Println("Sent data:", 45)
	}()

	// 主 goroutine 从 channel 接收数据
	for i := 0; i < 4; i++ {
		data := <-ch
		fmt.Println("Received data:", data)
	}
}
