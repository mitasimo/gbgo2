package main

import (
	"fmt"
	"os"
	"os/signal"
	"time"
)

func main() {

	// канал для завершения горутин
	done := make(chan struct{})

	for i := 0; i < 10; i++ {
		go func(num int, sleep time.Duration) {
			for {
				select {
				case <-done: // закрыт канал завершения
					fmt.Println("Завершилась горутина ", num)
					return // завершит горутину
				default:
					time.Sleep(sleep) // заснуть
				}
			}

		}(10-i, time.Millisecond*time.Duration(i+500))
	}

	// канал для сигналов ОС
	sig := make(chan os.Signal)
	// подписаться на сигнал SIGTERM
	signal.Notify(sig, os.Interrupt)

	<-sig // ожидать получения сигнала ОС

	close(done)             // закрыть канал для заверешения
	time.Sleep(time.Second) // подождать секунду
	os.Exit(-1)             // завершить приложение... можно было не вызывать...

}
