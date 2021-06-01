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

	for i := 0; i < 10; i++ {
		go Do(ctx, time.Millisecond*time.Duration((i+1)*500), fmt.Sprintf("Горутина %d завершена", i))
	}

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
