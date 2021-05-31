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
	go func(c context.Context) {
		for {
			select {
			case <-ctx.Done():
				// при закрытии канала сообщить о завершении горуиы
				fmt.Println("Горутина 1 завершена")
				return
			default:
				// заснуть на 2 секунды
				time.Sleep(time.Second * 15)
			}

		}
	}(ctx)

	// вторая горутина
	go func(c context.Context) {
		for {
			select {
			case <-ctx.Done():
				// при закрытии канала сообщить о завершении горуиы
				fmt.Println("Горутина 2 завершена")
				return
			default:
				// заснуть на 2 секунды
				time.Sleep(time.Millisecond * 100)
			}
		}
	}(ctx)

	// ожидать сигнал ОС
	<-chanSig
	// завершить контекст
	cancel()
	// подождать завершения горутин 1 секунда
	time.Sleep(time.Second)
	// завершить приложеие
	os.Exit(-1)
}
