package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	Terminate os.Signal = syscall.SIGTERM
)

func main() {

	// создать канал для получиния сигналов ОС
	chanSig := make(chan os.Signal)
	// настроить канал на получение сигналов SIGINT и SIGTERM
	signal.Notify(chanSig, os.Interrupt, Terminate)

	// создать контекст для завершения грутин
	ctx, cancel := context.WithCancel(context.Background())

	// первая горутина
	go Do(ctx, time.Second*15, "Горутина 1 завершена")
	// вторая горутина
	go Do(ctx, time.Microsecond*100, "Горутина 2 завершена")

	// ожидать сигнал ОС
	<-chanSig
	// завершить контекст
	cancel()
	// подождать завершения горутин 1 секунда
	time.Sleep(time.Second)
	// завершить приложеие
	os.Exit(-1)
}

// Do проверяет закрытие канала ctx.
// Если закрыт, выводит сообщение message.
// Если не закрыт, то засыпает на timeToSleep
func Do(ctx context.Context, timeToSleep time.Duration, message string) {
	for {
		select {
		case <-ctx.Done():
			// при закрытии канала сообщить о завершении горуиы
			fmt.Println(message)
			return
		default:
			// заснуть на очень короткий промежуток времени
			time.Sleep(timeToSleep)
		}
	}
}
