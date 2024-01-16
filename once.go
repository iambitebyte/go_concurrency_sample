package main

import (
	"fmt"
	"sync"
)

func main() {
	var once sync.Once

	// 函数 f 只会被执行一次
	for i := 0; i < 5; i++ {
		once.Do(func() {
			fmt.Println("Doing something")
		})
	}
}
