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
