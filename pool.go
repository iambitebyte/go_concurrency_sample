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
