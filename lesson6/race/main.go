// Курс Go. Уровень 2
// Урок 6
// Задание 3
// Смоделировать ситуацию “гонки”, и проверить программу на наличии “гонки”

package main

import (
	"fmt"
	"sync"
)

func main() {
	var (
		wg  sync.WaitGroup
		res = 100 // изменяемый ресурс
	)

	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			switch {
			case res < 100: // чтение значения
				res += 100 // установка значения
			default:
				res %= 7 // установка значения
			}
		}()
	}

	// ожидать завершение горутин
	wg.Wait()
	fmt.Println(res)
}
