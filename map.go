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
