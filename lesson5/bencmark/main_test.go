package main

import (
	"sync"
	"testing"
)

func reader(locker sync.Locker, wg *sync.WaitGroup, numIter int) {
	defer wg.Done()
	for i := 0; i < numIter; i++ {
		locker.Lock()
		locker.Unlock()
	}
}

func writer(locker sync.Locker, wg *sync.WaitGroup, numIter int) {
	defer wg.Done()
	for i := 0; i < numIter; i++ {
		locker.Lock()
		locker.Unlock()
	}
}

func BenchmarkMutexW10R90(b *testing.B) {
	var (
		wg sync.WaitGroup
		me sync.Mutex
	)

	// 1 писатель
	wg.Add(1)
	go writer(&me, &wg, 900)

	// 9 читателей
	for i := 0; i < 9; i++ {
		wg.Add(1)
		go reader(&me, &wg, 100)
	}

	wg.Wait()
}

func BenchmarkRWMutexW10R90(b *testing.B) {
	var (
		wg sync.WaitGroup
		me sync.RWMutex
	)

	// 1 писатель
	wg.Add(1)
	go writer(&me, &wg, 900)

	// 9 читателей
	for i := 0; i < 9; i++ {
		wg.Add(1)
		go reader(me.RLocker(), &wg, 100)
	}

	wg.Wait()
}

func BenchmarkMutexW50R50(b *testing.B) {
	var (
		wg sync.WaitGroup
		me sync.Mutex
	)

	// 5 писателей
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go writer(&me, &wg, 100)
	}

	// 5 читателей
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go reader(&me, &wg, 100)
	}

	wg.Wait()
}
func BenchmarkRWMutexW50R50(b *testing.B) {
	var (
		wg sync.WaitGroup
		me sync.RWMutex
	)

	// 5 писателей
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go writer(&me, &wg, 100)
	}

	// 5 читателей
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go reader(me.RLocker(), &wg, 100)
	}

	wg.Wait()
}

func BenchmarkMutexW90R10(b *testing.B) {
	var (
		wg sync.WaitGroup
		me sync.Mutex
	)

	// 9 писателей
	for i := 0; i < 9; i++ {
		wg.Add(1)
		go writer(&me, &wg, 100)
	}

	// 1 читатель
	wg.Add(1)
	go reader(&me, &wg, 900)

	wg.Wait()
}
func BenchmarkRWMutexW90R10(b *testing.B) {
	var (
		wg sync.WaitGroup
		me sync.RWMutex
	)

	// 9 писателей
	for i := 0; i < 9; i++ {
		wg.Add(1)
		go writer(&me, &wg, 100)
	}

	// 1 читатель
	wg.Add(1)
	go reader(me.RLocker(), &wg, 900)

	wg.Wait()
}
