// Курс GO. Уровень 2
// Домашнее задание к уроку 1

package main

import (
	"fmt"
)

func main() {
	var err error

	Task1()

	err = Task2()
	if err != nil {
		fmt.Println("Task 2:", err)
	}

	err = Task3()
	if err != nil {
		fmt.Println("Task 3:", err)
	}

	// пример 4
	err = Task4()
	if err != nil {
		fmt.Println("Task 4:", err)
	}
}
