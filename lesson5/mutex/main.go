package main

import (
	"fmt"
	"sync"
)

// SyncInt позволяет безопастно увеличивать свое значение из разных горутин
type SyncInt struct {
	sync.Mutex
	val int
}

// Add добавляет значение
func (si *SyncInt) Add(v int) int {
	// увеличить на v
	si.val += v
	// вернуть результат
	return si.val
}

// Val возвращает внутреннее значение
func (si *SyncInt) Val() int {
	return si.val
}

func main() {

	var (
		wg sync.WaitGroup // для ожидание завершения всех горутин
		si = SyncInt{}    // горутино безопастный int
	)

	for i := 1; i < 10; i++ {
		wg.Add(1) // увлеличить счетчик
		go func(si *SyncInt, inc, num int) {
			// после завершения функции уменьшить счетчик
			defer wg.Done()
			// заблокировать SyncInt, чтобы далее использовать безопастно
			si.Lock()
			// разблокировать в конце функции
			defer si.Unlock()

			olVal := si.Val()
			newVal := si.Add(inc)
			fmt.Printf("Горутина %d; Старое значение = %2d, инкремент = %2d, новое значение = %2d\n", num, olVal, inc, newVal)
		}(&si, 2*i+1, i)
	}

	// ожидать завершение всех горутин
	wg.Wait()
	// вывести результат
	// si.Lock()
	fmt.Printf("Значение = %d\n", si.Val())

}
