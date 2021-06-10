package main

import (
	"os"
	"runtime"
	"runtime/trace"
	"sync"
)

func main() {
	var (
		wg sync.WaitGroup
	)

	// начать запись трассировки
	trace.Start(os.Stderr)
	defer trace.Stop()

	for i := 0; i < 16; i++ {
		// увеличить счетчик
		wg.Add(1)
		go func() {
			// при завершении функции уменьшить счетчик
			defer wg.Done()

			for j := 0; j < 1000000; j++ {
				if j%100 == 0 { // на каждом 100м элементе
					// попросить планировщик прекратить выполнение потока и проверить,
					// нет ли других потоков в состоянии готовности
					runtime.Gosched()
				}
			}
		}()
	}

	// дождасться завершения всех горутин
	wg.Wait()
}
