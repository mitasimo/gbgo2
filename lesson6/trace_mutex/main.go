package main

import (
	"os"
	"runtime/trace"
	"sync"
)

func main() {
	var (
		wg sync.WaitGroup
		me sync.Mutex
	)

	// начать запись трассировки
	trace.Start(os.Stderr)
	defer trace.Stop()

	// мапа для совместного доступа из горутин
	mp := make(map[int]int)

	for i := 0; i < 16; i++ {
		wg.Add(1) // увеличить счетчик
		go func(num int) {
			defer wg.Done() // по завершении уменьшить счетчик
			for j := 0; j < 1000000; j++ {
				me.Lock()     // начало критической секции
				mp[num] = num // изменить мапу
				me.Unlock()   // конец критической секции
			}
		}(i)
	}

	wg.Wait() // дождаться завершения всех горутин

}
