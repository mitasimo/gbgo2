package main

import (
	"fmt"
	"time"
)

// Обертка над вызовом divideByZero с
//	- перехватом паники
//	- возвратом ошибки
func task2() (err error) {
	defer func() {
		// перехватить панику
		if v := recover(); v != nil {
			// обернем ошибку, которую вернет NewErrorWithTime()
			err = fmt.Errorf("%s; %w", v, NewErrorWithTime())
		}
	}()
	_ = divideByZero() // вызывать панику деления на 0
	return
}

// NewErrorWithTime - создает новую ошибку ErrorWithTime
func NewErrorWithTime() *ErrorWithTime {
	return &ErrorWithTime{Time: time.Now()}
}

// Тип ошибки, хранящий время ее возникновения
type ErrorWithTime struct {
	time.Time
}

// Реализация функции Error из интерфейса error
func (e *ErrorWithTime) Error() string {
	return fmt.Sprintf("an error occurred %s", e.Format("2006-01-02 15:04:05"))
}
