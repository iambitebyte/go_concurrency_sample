package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	// 创建一个无缓冲的 channel
	ch := make(chan struct{})

	// 启动两个 goroutine，其中一个等待 ch，另一个向 ch 发送信号
	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Goroutine 1: Waiting for singal")
		<-ch // 等待接收信号
		fmt.Println("Goroutine 1: Received signal")
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		fmt.Println("Goroutine 2: Waiting 1 second")
		time.Sleep(time.Second)
		fmt.Println("Goroutine 2: Sending signal")
		close(ch) // 关闭 channel，发送信号
	}()

	// 等待两个 goroutine 完成
	wg.Wait()
}
