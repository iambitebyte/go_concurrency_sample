package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

func main() {
	var wg sync.WaitGroup

	// 创建一个带缓冲的 channel
	ch := make(chan int, 3)

	// 启动多个 goroutine 从 channel 接收数据
	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			for {
				data, ok := <-ch
				if !ok {
					fmt.Printf("Goroutine %s: Channel closed\n", id)
					return
				}
				fmt.Printf("Goroutine %s: Received data %d\n", id, data)
				time.Sleep(time.Second)
			}
		}(strconv.Itoa(i))
	}

	// 向 channel 发送数据
	for i := 1; i <= 5; i++ {
		ch <- i
	}

	// 关闭 channel，通知所有 goroutine 停止接收
	close(ch)

	// 等待所有 goroutine 完成
	wg.Wait()
}
