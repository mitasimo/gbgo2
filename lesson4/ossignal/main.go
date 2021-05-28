package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	// создать контекст, который закроет свой канал done() при поступлении сигнала SIGTERM
	ctxSig, cancelSig := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
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
