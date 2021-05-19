package main

import (
	"fmt"
	"os"
)

// task3 пытается создать 1 миллион пустых файлов
// конкретно на моем компьютере создается 8165 файлов
func task3() (err error) {
	const (
		numFiles = 1000000         // количество файлов
		dirName  = "../onemillion" // каталог для файлов
	)

	// перехватить все паники.. (вызванные в ветке вызовов этой функции)
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("an error occured: %v", v)
		}
	}()

	// создать каталог для файлов
	err = os.Mkdir(dirName, os.ModeDir|os.ModePerm)
	if err != nil {
		return
	}

	// цикл создания файлов
	for i := 0; i < numFiles; i++ {

		// сформировать имя файла
		fileName := fmt.Sprintf("%s/f_%07d", dirName, i)

		// создать файл
		f, errIn := os.Create(fileName)
		if errIn != nil {
			return errIn
		}

		// инициировать отложенное закрытие файла
		defer f.Close()
	}

	return
}
