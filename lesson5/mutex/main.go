package main

import (
	"fmt"
	"sync"
)

// SyncInt позволяет безопастно увеличивать свое значение из разных горутин
type SyncInt struct {
	val int
	l   sync.Mutex
}

// Add добавляет значение
func (si *SyncInt) Add(v int) int {
	// установить блокировку
	si.l.Lock()
	// после выполнения функции снять блокировку
	defer si.l.Unlock()
	// увеличить на v
	si.val += v
	// вернуть результат
	return si.val
}

// Val возвращает внутреннее значение
func (si *SyncInt) Val() int {
	si.l.Lock()
	defer si.l.Unlock()
	return si.val
}

func main() {

	var (
		wg sync.WaitGroup // для ожидание завершения всех горутин
		si = &SyncInt{}   // горутино безопастный int
	)

	for i := 1; i < 100; i++ {
		wg.Add(1) // увлеличить счетчик
		go func(si *SyncInt, val, num int) {
			// после завершения функции уменьшить счетчик
			defer wg.Done()
			// выполнить горутино безопастный Add
			fmt.Printf("Горутина %d: Значение: %d\n", num, si.Add(val))
		}(si, 2*i+1, i)
	}

	// ожидать завершение всех горутин
	wg.Wait()
	// вывести результат
	fmt.Printf("Значение = %d", si.Val())

}
