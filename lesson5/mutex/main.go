package main

import (
	"fmt"
	"sync"
)

// SyncInt позволяет безопастно увеличивать свое значение из разных горутин
type SyncInt struct {
	val  int
	lock sync.RWMutex
}

// Add добавляет значение
func (si *SyncInt) Add(v int) int {
	// установить блокировку
	si.lock.Lock()
	// после выполнения функции снять блокировку
	defer si.lock.Unlock()
	// увеличить на v
	si.val += v
	// вернуть результат
	return si.val
}

// Val возвращает внутреннее значение
func (si *SyncInt) Val() int {
	si.lock.RLock()
	defer si.lock.RUnlock()

	return si.val
}

func main() {

	var (
		wg sync.WaitGroup // для ожидание завершения всех горутин
		si = SyncInt{}    // горутино безопастный int
	)

	for i := 1; i < 10; i++ {
		wg.Add(1) // увлеличить счетчик
		go func(si *SyncInt, val, num int) {
			// после завершения функции уменьшить счетчик
			defer wg.Done()
			// выполнить горутино безопастный Add
			newVal := si.Add(val)
			fmt.Printf("Горутина %d изменила значение на %3d. Новое значение = %2d\n", num, val, newVal)
		}(&si, 2*i+1, i) // 2*i+1
	}

	// ожидать завершение всех горутин
	wg.Wait()
	// вывести результат
	fmt.Printf("Значение = %d\n", si.Val())

}
