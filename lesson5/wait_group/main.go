// GO. Уровень 2
// Урок 5
// Задание 1
// Напишите программу, которая запускает n потоков и дожидается завершения их всех

package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {

	var wg sync.WaitGroup
	for i := 0; i < 10; i += 2 {
		go Do(&wg, 20*i, i)
	}

	wg.Wait()
	fmt.Println("Программа завершена")

}

func Do(wg *sync.WaitGroup, millisecondToSleep, num int) {
	wg.Add(1)
	defer wg.Done()

	fmt.Printf("Выполняется горутина %2d\n", num)
	time.Sleep(time.Millisecond * time.Duration(millisecondToSleep))
}
