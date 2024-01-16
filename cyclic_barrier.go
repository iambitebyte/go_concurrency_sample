package main

import (
	"fmt"
	"sync"
)

type CyclicBarrier struct {
	count int
	total int
	cond  *sync.Cond
	mut   sync.Mutex
}

func NewCyclicBarrier(total int) *CyclicBarrier {
	b := &CyclicBarrier{
		count: 0,
		total: total,
		cond:  sync.NewCond(&sync.Mutex{}),
		mut:   sync.Mutex{},
	}
	return b
}

func (b *CyclicBarrier) Await() {
	b.mut.Lock()
	defer b.mut.Unlock()

	b.count++
	if b.count < b.total {
		b.cond.Wait()
	} else {
		// Reset the barrier
		b.count = 0
		b.cond.Broadcast()
	}
}

func main() {
	const totalThreads = 3
	barrier := NewCyclicBarrier(totalThreads)
	var wg sync.WaitGroup

	for i := 0; i < totalThreads; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			fmt.Printf("Goroutine %d: Waiting at the barrier\n", id)
			barrier.Await()
			fmt.Printf("Goroutine %d: Passed the barrier\n", id)
		}(i)
	}

	wg.Wait()
}
