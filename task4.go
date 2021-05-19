package main

import (
	"fmt"
	"time"
)

func task4() (err error) {
	// этим кодом паника обработана не будет!!!!
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("%v", v)
		}
	}()

	go func() {
		// recover должен вызываться в той же горутине, в которой вызывается паника!!!!
		defer func() {
			if v := recover(); v != nil {
				// err - это возвращаемое значение функции task4()
				err = fmt.Errorf("%v", v)
			}
		}()
		// паникуем в горутине!!!
		panic("Panic in gorouting!!!")
	}()

	// подождать завершения горутины
	time.Sleep(time.Second * 2)

	return
}
