// GO. Уровень 2
// Урок 5
// Задание 3
//
// файл содержит бенчмарки
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

// Bench запускает выполнение теста
// 	name - имя теста
//	b -
//  s - множество
//	numReaders - количество читателей
//	numWriters - количество писателей
//
func Bench(name string, b *testing.B, s Set, numReaders, numWriters int) {
	// Магическая константа!!!
	// testing.B.Run() запускает переданный ей func() 5 раз
	const numRun = 5

	// всего читателей и писателей
	numTotal := numReaders + numWriters
	// всего задач
	numTasks := runtime.GOMAXPROCS(0) * numTotal * numRun
	// буферированный канал для получения задач в тесте
	fnChan := make(chan func(), numTasks)

	// добавить в канал задачи читатели
	for i := 0; i < numTasks/numTotal*numReaders; i++ {
		fnChan <- makeReader(s)
	}

	// добавитьв канал задачи писатели
	for i := 0; i < numTasks/numTotal*numWriters; i++ {
		fnChan <- makeWriter(s)
	}

	b.ResetTimer()
	b.Run(name, func(b1 *testing.B) {
		// установить уровень параллелизма
		b1.SetParallelism(numTotal)
		// запустить параллельные задачи
		b1.RunParallel(func(pb *testing.PB) {
			// получить задачу из канала
			fn := <-fnChan
			for pb.Next() {
				// выполнить задачу
				fn()
			}
		})
	})

	// закрыть канал
	close(fnChan)
}

func BenchmarkSet(b *testing.B) {

	setMutex := NewSetMutex()
	setRWMutex := NewSetRWMutex()
	setAtomic := NewSetAtomic()

	// по факту читателей и писателей будет больше в runtime.GOMAXPROCS(0) * 5 раз
	// Параметры Bench() numReaders и numWriters задают по сути процентное отношение

	Bench("SetMutex: 1 читатель, 9 писателей", b, setMutex, 1, 9)
	Bench("SetRWMutex: 1 читатель, 9 писателей", b, setRWMutex, 1, 9)
	Bench("SetAtomic: 1 читатель, 9 писателей", b, setAtomic, 1, 9)

	Bench("SetMutex: 5 читателей, 5 писателей", b, setMutex, 5, 5)
	Bench("SetRWMutex: 5 читателей, 5 писателей", b, setRWMutex, 5, 5)
	Bench("SetAtomic: 5 читателей, 5 писателей", b, setAtomic, 5, 5)

	Bench("SetMutex: 9 читателей, 1 писатель", b, setMutex, 9, 1)
	Bench("SetRWMutex: 9 читателей, 1 писатель", b, setRWMutex, 9, 1)
	Bench("SetAtomic: 9 читатель, 1 писателей", b, setAtomic, 9, 1)

}
