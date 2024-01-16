package main

import (
	"fmt"
	"sync"
	"time"
	"golang.org/x/sync/singleflight"
)

func main() {
	var sf singleflight.Group
	var wg sync.WaitGroup

	// 启动多个 goroutine 同时请求相同的数据
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()

			// 请求相同的数据
			val, err, _ := sf.Do("key", func() (interface{}, error) {
				// 模拟耗时操作
				time.Sleep(100 * time.Millisecond)
				return id, nil
			})

			if err != nil {
				fmt.Printf("Goroutine %d: Error: %v\n", id, err)
			} else {
				fmt.Printf("Goroutine %d: Result: %v\n", id, val)
			}
		}(i)
	}

	wg.Wait()
}
