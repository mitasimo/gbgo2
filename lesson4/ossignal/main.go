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
				time.Sleep(time.Second * 2)
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
				time.Sleep(time.Second * 2)
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

	// создать контекст, который закроет свой канал done() при поступлении сигнал SIGINT и SIGTERM
	ctxSig, cancelSig := signal.NotifyContext(context.Background(), os.Interrupt, Terminate)
	// в конце функции main вызвать функцию завершения
	defer cancelSig()

	// создать контектс с таймаутом
	ctxTimeoOut, cancelTimeoOut := context.WithTimeout(context.Background(), time.Second)
	// в конце функции main вызвать функцию завершения
	defer cancelTimeoOut()

	// ожидать поступления сигнала SIGTERM или завершения таймаута
	select {
	case <-ctxSig.Done():
		fmt.Printf("Поступил сигнал. Причина завершения контекста: %v\n", ctxSig.Err())
	case <-ctxTimeoOut.Done():
		fmt.Printf("Наступил таймаут. Причина завершения контекста:  %v\n", ctxTimeoOut.Err())
	}
}
