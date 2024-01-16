# go_concurrency_sample

## Goroutine
Go语言的协程（Goroutine）是一种轻量级线程，由 Go 运行时系统（Runtime）管理。

对于内核来说，协程们会串行的运行于单个线程，类似node.js的异步I/O与事件驱动的机制。
go的协程调度完全是由go运行时在`用户态`上完成的，因此在高并发的场景下，它比多线程技术要更加高效。协程的调度不需要内核参与，省去了内核与用户态的切换，一定程度上更适合执行过程短，但并发量高的场景。


在`goroutine.go`中，启动协程的方法是`go {函数名}`
```
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

```

## RWMutex
`sync.RWMutex`是Go语言标准库中提供的读写锁（Read-Write Mutex），用于在多个 goroutine 之间提供对共享资源的安全访问。读写锁分为读锁和写锁，多个 goroutine 可以同时持有读锁，但只有一个 goroutine 可以持有写锁。

读锁可以同时持有，不过读锁也写锁之间也是互斥的。

RWMutex 是一种强大的同步机制，可以在读取频繁的情况下提高并发性能，因为多个 goroutine 可以同时获取读锁。然而，必须小心在使用 RWMutex 时防止死锁。
```
package main

import (
	"fmt"
	"sync"
	"time"
)

var (
	data  map[string]string
	mutex sync.RWMutex
)

func main() {
	data = make(map[string]string)

	// go writeData("key1", "a", "go")
	// go readData("key1", "a", 0)
	// go readData("key1", "b", 1)
	// go writeData("key1", "b", "stop")

	go writeData("key1", "a", "kill bill.")
	go readData("key1", "a", 0)
	go readData("key1", "b", 0)

	time.Sleep(5 * time.Second)
}

func readData(key string, name string, second time.Duration) {
	mutex.RLock()
	fmt.Printf("%s(%s): Attempting to read data\n", name, time.Now().String())
	defer mutex.RUnlock()
	time.Sleep(second * time.Second)
	value := data[key]
	fmt.Printf("%s(%s): Read data: %s=%s\n", name, time.Now().String(), key, value)
}

func writeData(key, name string, value string) {

	mutex.Lock()
	fmt.Printf("%s(%s): Attempting to write data\n", name, time.Now().String())
	defer mutex.Unlock()

	data[key] = value
	time.Sleep(time.Second)
	fmt.Printf("%s(%s): Write data: %s=%s\n", name, time.Now().String(), key, value)
}

```

## Cond

在Go语言中，`sync.Cond（条件变量）`是一种在多个 goroutine 之间进行协调的机制。sync.Cond 用于等待或唤醒 goroutine，通常与互斥锁（sync.Mutex）一起使用，以在共享资源的访问上进行同步。 

在`cond.go`的例子中，一旦某个协程调用`cond.wait()`，就会将运行的机会让给其他的协程。如果所有协程
```
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

	// 主 goroutine 等待一段时间，然后唤醒所有等待的 goroutine
	time.Sleep(2 * time.Second)
	cond.Signal() // 唤醒一个等待的 goroutine
	time.Sleep(2 * time.Second)
	cond.Broadcast() // 唤醒所有等待的 goroutine

	time.Sleep(5 * time.Second)
}

func producer(name string, delay time.Duration) {
	for {
		mutex.Lock()
		fmt.Printf("%s: Waiting for condition\n", name)
		cond.Wait() // 阻塞等待条件变为真
		fmt.Printf("%s: Condition satisfied, doing work\n", name)
		mutex.Unlock()

		time.Sleep(delay * time.Second)
	}
}
```

## Once

在 Go 语言中，`sync.Once` 是一个同步工具，用于确保某个操作只执行一次。sync.Once 提供了一种机制，使得在并发程序中某个函数只执行一次，无论有多少个 goroutine 调用它。

sync.Once 使用了一个内部的 done 变量来跟踪函数是否已经执行。当 Do 方法被调用时，它会检查 done 变量，如果尚未执行过，就执行传入的函数，然后将 done 设置为已执行。在之后的调用中，Do 方法会直接返回，不再执行传入的函数。

在这个例子中，无论 Do 方法被调用多少次，传入的函数只会执行一次。输出结果只有一行 "Doing something"。

sync.Once 主要用于初始化操作，确保某个初始化函数只被调用一次。这在全局或包级别的变量初始化中特别有用。例如，你可以在全局变量的初始化函数中使用 sync.Once 确保初始化只执行一次，即使有多个 goroutine 同时尝试初始化。

```
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
```

## map
在 Go 语言中，map 是一种内置的数据结构，用于存储键值对。map 提供了一种方便且高效的方式来表示和操作无序的集合。它类似于其他编程语言中的字典或关联数组。
以下是 map 的主要特点和用途：

1. 键值对存储：map 可以存储键值对，其中每个键都必须是唯一的，而值则可以重复。
2. 快速查找： 使用键作为索引，可以快速查找对应的值。这使得 map 在查找和检索数据时非常高效。
3. 动态增长：map 可以动态增长，不需要提前声明容量。在运行时，可以方便地添加或删除键值对。
4. 无序性：map 是无序的，即插入顺序不保证与遍历顺序一致。如果需要有序的键值对集合，可以考虑使用 slice 配合排序。


map 在 Go 中是一个非常常用的数据结构，用于存储和操作键值对数据。在处理配置信息、缓存、数据索引等场景中，map 都是一种非常实用的选择。但需要注意，map 并不是并发安全的，如果在多个 goroutine 中同时读写同一个 map，需要采取额外的同步措施，比如使用互斥锁（sync.Mutex）或并发安全的 sync.Map。  

```
package main

import "fmt"

func main() {
	// 创建一个空的 map
	myMap := make(map[string]int)

	// 添加键值对
	myMap["one"] = 1
	myMap["two"] = 2
	myMap["three"] = 3

	// 获取值
	fmt.Println("Value for key 'two':", myMap["two"])

	// 遍历 map
	for key, value := range myMap {
		fmt.Printf("Key: %s, Value: %d\n", key, value)
	}

	// 删除键值对
	delete(myMap, "two")

	// 检查键是否存在
	if value, exists := myMap["two"]; exists {
		fmt.Println("Value for key 'two':", value)
	} else {
		fmt.Println("Key 'two' not found.")
	}
}

```

## Pool

在 Go 语言中，sync.Pool 是用于对象池的同步工具。它提供了一个简单的、线程安全的机制，用于重用对象，以减轻垃圾收集的压力，并提高性能。

sync.Pool 主要有两个方法：Put 用于向池中存放对象，Get 用于从池中获取对象。池的特点是在高并发的情况下，能够有效地重用对象，而不是频繁地创建和销毁对象。

在这个例子中，sync.Pool 用于存放 MyObject 对象。当调用 Get 方法时，它会尝试从池中获取一个对象。如果池中没有可用的对象，它会调用 New 函数生成一个新的对象。当调用 Put 方法时，它将对象放回池中，以供下一次调用 Get 时重用。

使用对象池的好处在于可以避免频繁地创建和销毁对象，特别是在高并发的情况下，可以显著提高性能。需要注意的是，对象池中的对象并没有固定的生命周期，垃圾收集器会在适当的时机回收这些对象。

```
package main

import (
	"fmt"
	"sync"
)

type MyObject struct {
	Data int
}

func main() {
	// 创建一个对象池
	objectPool := &sync.Pool{
		New: func() interface{} {
			// 当池为空时，New 函数会被调用来生成新的对象
			fmt.Println("Creating a new object")
			return &MyObject{}
		},
	}

	// 从池中获取对象
	obj1 := objectPool.Get().(*MyObject)
	obj1.Data = 42
	fmt.Printf("Object 1: %+v\n", obj1)

	// 将对象放回池中
	objectPool.Put(obj1)

	// 再次从池中获取对象，这时应该是之前放回去的对象
	obj2 := objectPool.Get().(*MyObject)
	fmt.Printf("Object 2: %+v\n", obj2)
	objectPool.Put(obj2)

	//放回后连续获得对象
	obj3 := objectPool.Get().(*MyObject)
	fmt.Printf("Object 3: %+v\n", obj3)

	obj4 := objectPool.Get().(*MyObject)
	fmt.Printf("Object 4: %+v\n", obj4)
}
```

## Context

在 Go 语言中，context 是一个标准库中的包，用于在 goroutine 之间传递取消信号、截止时间以及其他请求范围的值。context 包提供了一个 Context 接口和一些用于创建和处理 Context 的函数。

```
package main

import (
	"context"
	"fmt"
	"time"
)

func main() {
	// 创建一个带有截止时间的 Context
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel() // 确保在函数退出时取消

	// 启动一个 goroutine，模拟耗时操作
	go func() {
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
```

## channel

在 Go 语言中，channel 是一种并发编程的原语，用于在不同的 goroutine 之间安全地传递数据和同步执行。channel 提供了一种通信的方式，可以让多个 goroutine 之间进行数据交换，避免了显式的锁和共享内存，使得编写并发程序更加简洁和安全。  

- `数据传递(channel1.go)` : channel可以用于在不同的 goroutine 之间传递数据。一个 goroutine 可以将数据发送到 channel，而另一个 goroutine 可以从 channel 中接收到这些数据。
- `同步执行(channel2.go)` : channel 也可以用于在多个 goroutine 之间进行同步，确保某个操作发生在另一个操作之前。通过 channel 的阻塞特性，可以实现简单的同步。
- `关闭和广播(channel3.go)` : 通过关闭 channel 可以向多个 goroutine 发送广播，告诉它们停止工作。当一个 channel 被关闭后，任何尝试向其发送数据的操作都会立即完成，而接收操作会返回已关闭的 channel 的零值。 

## Semaphore

在 Go 语言中，通常使用 sync 包中的 Semaphore（信号量）来控制并发访问的数量。虽然 sync 包本身并没有提供显式的信号量类型，但可以使用 chan struct{} 来模拟信号量的行为，实现类似信号量的效果。  

在这个例子中，semaphore 是一个带有容量的通道，通道的容量表示可以同时执行的任务数量。通过向通道发送 struct{} 来获取信号量，通过从通道接收 struct{} 来释放信号量。

这种模拟信号量的方式可以帮助控制并发访问的数量，防止过多的 goroutine 同时执行，以确保资源的合理利用。在实际的应用中，可以根据具体需求调整信号量的容量和控制并发的逻辑。

```
package main

import (
	"fmt"
	"sync"
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
			fmt.Printf("Goroutine %d: Running\n", id)
		}(i)
	}

	wg.Wait()
}
```

## SingleFlight

在 Go 语言中，Singleflight 是一个用于合并相同请求的库，它可以有效地减轻并发场景下的重复请求。Singleflight 通过在多个 goroutine 之间共享相同请求的结果，以降低并发负载。这个库是由 Uber 提供的，并且可以在 GitHub 上找到：github.com/uber-go/singleflight。

Singleflight 主要通过在调用某个函数时检查是否已经有其他 goroutine 在执行相同的函数，如果是，则等待该函数的结果，而不是重新执行。这样可以避免多个 goroutine 同时执行相同的昂贵操作，从而减轻负载。

在上述例子中，通过 sf.Do 函数，所有请求相同的数据的 goroutine 会共享相同的结果，避免了重复执行相同的耗时操作。Singleflight 在某些场景下可以帮助提高系统的并发性能，减轻负载。  

运行`single_flight.go`需要先安装第三方模块

`go get golang.org/x/sync/singleflight`

```
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
```

## CyclicBarrier

在 Go 语言中，并没有内置类似于 Java 中的 CyclicBarrier工具。CyclicBarrier 是一种同步工具，用于在多个线程之间建立屏障，当所有线程都达到屏障时，屏障会打开，所有线程可以继续执行。

尽管 Go 没有直接提供 CyclicBarrier，但你可以使用其他 Go 语言提供的同步工具来实现类似的功能，例如使用 sync.WaitGroup 和 sync.Cond。

```
package main

import (
	"fmt"
	"sync"
)

type CyclicBarrier struct {
	count  int
	total  int
	cond   *sync.Cond
	mut    sync.Mutex
}

func NewCyclicBarrier(total int) *CyclicBarrier {
	b := &CyclicBarrier{
		count:  0,
		total:  total,
		cond:   sync.NewCond(&sync.Mutex{}),
		mut:    sync.Mutex{},
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

```