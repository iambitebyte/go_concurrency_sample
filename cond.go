package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	mutex sync.Mutex
	cond  *sync.Cond
)

func main() {
	cond = sync.NewCond(&mutex)

	go producer("A", 2)
	go producer("B", 3)

	time.Sleep(2 * time.Second)

	// cond.Signal()
	// cond.Broadcast()

	time.Sleep(2 * time.Second)

	cond.Broadcast()

	time.Sleep(5 * time.Second)
}

func producer(name string, delay time.Duration) {
	mutex.Lock()
	fmt.Printf("%s-%s: waiting for condition\n", name, time.Now().String())

	time.Sleep(time.Second * 2)

	cond.Wait()
	fmt.Printf("%s-%s: Condition satisfied, doing work\n", name, time.Now().String())
	mutex.Unlock()
	time.Sleep(delay * time.Second)
	fmt.Printf("%s-%s: Delay finshed\n", name, time.Now().String())
}
