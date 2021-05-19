package main

import (
	"fmt"
)

// divideByZero - генерирует неявную ошибку (деление на 0)
func divideByZero() int {
	var zero int
	return 1 / zero
}

// Обертка над вызовом divideByZero
func task1() {
	defer func() {
		// перехватить панику
		if v := recover(); v != nil {
			// вывести ошибку в консоль
			fmt.Printf("Task 1: %v\n", v)
		}
	}()
	_ = divideByZero() // вызывать панику деления на 0
}
