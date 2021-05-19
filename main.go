package main

import (
	"fmt"
)

func main() {
	var err error

	task1()

	err = task2()
	if err != nil {
		fmt.Println("Task 2:", err)
	}

	err = task3()
	if err != nil {
		fmt.Println("Task 3:", err)
	}

	// пример 4
	err = task4()
	if err != nil {
		fmt.Println("Task 4:", err)
	}
}
