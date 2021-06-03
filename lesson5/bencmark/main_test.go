package main

import (
	"sync"
	"testing"
)

func reader(resourse *int, locker sync.Locker, wg *sync.WaitGroup, numIter int) {
	var v int
	defer wg.Done()
	for i := 0; i < numIter; i++ {
		locker.Lock()
		v = *resourse
		v++
		locker.Unlock()
	}
}

func writer(resourse *int, locker sync.Locker, wg *sync.WaitGroup, numIter int) {
	defer wg.Done()
	for i := 0; i < numIter; i++ {
		locker.Lock()
		*resourse++
		locker.Unlock()
	}
}

func BenchmarkMutexW10R90(b *testing.B) {
	var (
		wg  sync.WaitGroup
		me  sync.Mutex
		res int
	)

	// 1 писатель
	wg.Add(1)
	go writer(&res, &me, &wg, 90)

	// 9 читателей
	for i := 0; i < 9; i++ {
		wg.Add(1)
		go reader(&res, &me, &wg, 10)
	}

	wg.Wait()
}

func BenchmarkRWMutexW10R90(b *testing.B) {
	var (
		wg  sync.WaitGroup
		me  sync.RWMutex
		res int
	)

	// 1 писатель
	wg.Add(1)
	go writer(&res, &me, &wg, 90)

	// 9 читателей
	for i := 0; i < 9; i++ {
		wg.Add(1)
		go reader(&res, me.RLocker(), &wg, 10)
	}

	wg.Wait()
}

func BenchmarkMutexW50R50(b *testing.B) {
	var (
		wg  sync.WaitGroup
		me  sync.Mutex
		res int
	)

	// 5 писателей
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go writer(&res, &me, &wg, 10)
	}

	// 5 читателей
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go reader(&res, &me, &wg, 10)
	}

	wg.Wait()
}
func BenchmarkRWMutexW50R50(b *testing.B) {
	var (
		wg  sync.WaitGroup
		me  sync.RWMutex
		res int
	)

	// 5 писателей
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go writer(&res, &me, &wg, 10)
	}

	// 5 читателей
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go reader(&res, me.RLocker(), &wg, 10)
	}

	wg.Wait()
}

func BenchmarkMutexW90R10(b *testing.B) {
	var (
		wg  sync.WaitGroup
		me  sync.Mutex
		res int
	)

	// 9 писателей
	for i := 0; i < 9; i++ {
		wg.Add(1)
		go writer(&res, &me, &wg, 10)
	}

	// 1 читатель
	wg.Add(1)
	go reader(&res, &me, &wg, 90)

	wg.Wait()
}
func BenchmarkRWMutexW90R10(b *testing.B) {
	var (
		wg  sync.WaitGroup
		me  sync.RWMutex
		res int
	)

	// 9 писателей
	for i := 0; i < 9; i++ {
		wg.Add(1)
		go writer(&res, &me, &wg, 10)
	}

	// 1 читатель
	wg.Add(1)
	go reader(&res, me.RLocker(), &wg, 90)

	wg.Wait()
}
