package main

import (
	"fmt"
	"time"
)

// Task2 - код для задания 2
func Task2() (err error) {
	defer func() {
		// перехватить панику
		if v := recover(); v != nil {
			// обернем ошибку, которую вернет NewErrorWithTime()
			err = fmt.Errorf("%s; %w", v, NewErrorWithTime())
		}
	}()
	_ = DivideByZero() // вызывать панику деления на 0
	return
}

// NewErrorWithTime создает новую ошибку ErrorWithTime
func NewErrorWithTime() *ErrorWithTime {
	return &ErrorWithTime{Time: time.Now()}
}

// ErrorWithTime - тип ошибки, хранящий время ее возникновения
type ErrorWithTime struct {
	time.Time
}

// Error - реализация функции Error из интерфейса error
func (e *ErrorWithTime) Error() string {
	return fmt.Sprintf("an error occurred %s", e.Format("2006-01-02 15:04:05"))
}
