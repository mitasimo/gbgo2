package main

import (
	"runtime"
	"testing"
)

// Set интерфейс определяющий операции над множеством
type Set interface {
	Add(int)
	Has(int) bool
}

// makeReader возвращает функцию чтения из множества
func makeReader(s Set) func() {
	return func() {
		s.Has(0)
	}
}

// makeWriter - возвращает функция записи числа в множество
func makeWriter(s Set) func() {
	return func() {
		s.Add(0)
	}
}

func Bench(name string, b *testing.B, s Set, numReaders, numWriters int) {
	// Магическая константа!!!
	// testing.B.Run() запускает переданный ей func() 5 раз
	const numRun = 5

	// всего задач
	numTotal := numReaders + numWriters

	numTasks := runtime.GOMAXPROCS(0) * numTotal * numRun

	fnChan := make(chan func(), numTasks)

	for i := 0; i < numTasks/numTotal*numReaders; i++ {
		fnChan <- makeReader(s)
	}

	for i := 0; i < numTasks/numTotal*numWriters; i++ {
		fnChan <- makeWriter(s)
	}

	b.ResetTimer()
	b.Run(name, func(b1 *testing.B) {
		// atomic.AddInt64(&cnt, 1)
		b1.SetParallelism(numTotal)
		b1.RunParallel(func(pb *testing.PB) {
			fn := <-fnChan
			for pb.Next() {
				fn()
			}
		})
	})

	close(fnChan)
}

func BenchmarkSet(b *testing.B) {

	setMutex := NewSetMutex()
	setRWMutex := NewSetRWMutex()

	Bench("SetMutex: 1 читатель, 9 писателей", b, setMutex, 1, 9)
	Bench("SetRWMutex: 1 читатель, 9 писателей", b, setRWMutex, 1, 9)

	Bench("SetMutex: 5 читателей, 5 писателей", b, setMutex, 5, 5)
	Bench("SetRWMutex: 5 читателей, 5 писателей", b, setRWMutex, 5, 5)

	Bench("SetMutex: 9 читателей, 1 писатель", b, setMutex, 9, 1)
	Bench("SetRWMutex: 9 читателей, 1 писатель", b, setRWMutex, 9, 1)

}
