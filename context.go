package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// 创建一个带有截止时间的 Context
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel() // 确保在函数退出时取消

	// 启动一个 goroutine，模拟耗时操作
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel() // 确保在函数退出时取消

		select {
		case <-ctx.Done():
			fmt.Println("Main goroutine:", ctx.Err())
		}

		// 模拟耗时操作，超过截止时间将被取消
		time.Sleep(3 * time.Second)
		fmt.Println("Goroutine finished")
	}()

	// 在主 goroutine 中等待 Context 的完成信号
	select {
	case <-ctx.Done():
		fmt.Println("Main goroutine:", ctx.Err())
	}
}
