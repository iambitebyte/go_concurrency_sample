package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	const concurrentLimit = 3

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, concurrentLimit)

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			semaphore <- struct{}{} // 获取信号量，类似于 P 操作
			defer func() {
				<-semaphore // 释放信号量，类似于 V 操作
			}()

			// 执行并发任务
			fmt.Printf("%s: Goroutine %d: Running\n", time.Now().String(), id)
			time.Sleep(time.Second)
		}(i)
	}

	wg.Wait()
}
