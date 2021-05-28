package main

import (
	"fmt"
)

// DivideByZero генерирует неявную ошибку (деление на 0)
func DivideByZero() int {
	var zero int
	return 1 / zero
}

// Task1 - код для задания 1
func Task1() {
	defer func() {
		// перехватить панику
		if v := recover(); v != nil {
			// вывести ошибку в консоль
			fmt.Printf("Task 1: %v\n", v)
		}
	}()
	_ = DivideByZero() // вызывать панику деления на 0
}
