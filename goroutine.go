package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	// 创建一个等待组，用于等待所有协程完成
	var wg sync.WaitGroup

	// 启动多个协程
	for i := 1; i <= 3; i++ {
		wg.Add(1) // 每启动一个协程，等待组计数加一
		go worker(i, &wg)
	}

	// 等待所有协程完成
	wg.Wait()

	fmt.Println("All Goroutines have finished.")
}

func worker(id int, wg *sync.WaitGroup) {
	// 在协程结束时通知等待组减一
	defer wg.Done()

	fmt.Printf("Goroutine %d started.\n", id)

	// 模拟协程执行任务
	time.Sleep(2 * time.Second)

	fmt.Printf("Goroutine %d finished.\n", id)
}
