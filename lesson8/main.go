// Утилита для поиска и удаления задублированных файлов
// в каталоге и его подкаталогах
// Вызывается
// $ util -p=<стартовый каталог> -r
// ключ -dir (directory) задает стартовый каталог (по уполчанию текущий каталог программы)
// ключ -rm (remove) заставит утилиту удалить файлы (оставит только один)

package main

import (
	"flag"
	"fmt"
	"hash/adler32"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
)

// FileHash связывает путь к файлу и его хэш
type FileHash struct {
	Hash     uint32
	FilePath string
}

func main() {

	var (
		startPath    string // каталог с которого начинать перебор файлов
		removeCopies bool   // признак необходимости удаления копий
	)

	// чтение флогов командной строки
	flag.StringVar(&startPath, "dir", ".", "начальный каталог")
	flag.BoolVar(&removeCopies, "rm", false, "удалять копии")
	flag.Parse()

	if startPath == "" {
		log.Fatal("не задан начальный каталог")
	}

	fmt.Println("Начальный каталог:", startPath)

	var (
		wg sync.WaitGroup
	)

	filePathChan := make(chan string)    // канал для передачи путей к файлам
	fileHashChan := make(chan *FileHash) // канал для передачи хеша файлов

	// запуск горутину=ы, считывающей имена файлав из каталога
	go func() {
		// после завершения итераций по файлам закрыть канал путей к файлам
		defer close(filePathChan)
		IterateEntitiesInDirectory(startPath, filePathChan)
	}()

	// запуск горутин подсчета хеша файлов
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for filePath := range filePathChan { // получить очередной путь к файлу из канала
				// открыть файл по пути
				file, err := os.Open(filePath)
				if err != nil {
					log.Println(err) // записать ошибку в лог
					continue         // пропустить путь к файлу
				}

				// хешировать данные файла
				hash := adler32.New()
				io.Copy(hash, file)
				file.Close()

				// отправить хеш файла в канал
				fileHashChan <- &FileHash{
					FilePath: filePath,
					Hash:     hash.Sum32(),
				}
			}
		}()
	}

	// Очень важная горутина!!!
	// Дожидается завершения горутин расчета хеша,
	// после чего закрывает канал для хешей.
	// Без нее основаная горутина всегда будет
	// ожидать хеш из канала
	go func() {
		wg.Wait()
		close(fileHashChan)
	}()

	// мапа для сбора данных о дублях
	// ключ - хеш uint32
	// значение - массив путей к одинаковым файлам
	copies := make(map[uint32][]string)

	// читаем из канала хеши с их путями и добавляем в мапу копий
	for fileHash := range fileHashChan {
		filesPath := copies[fileHash.Hash]               // получить массив путей к файлам
		filesPath = append(filesPath, fileHash.FilePath) // добавить путь к массиву
		copies[fileHash.Hash] = filesPath                // сохранить новый массив путей
	}

	copiesPrinted := 0 // количество напечатанных путей к дублям

	for key, pathes := range copies {
		if len(pathes) > 1 { // есть копии
			copiesPrinted++
			fmt.Println("Хеш:", key) // вывести хеш
			for _, curPath := range pathes {
				fmt.Println("\t", curPath) // вывести путь к файлу
			}
		}
	}

	if copiesPrinted == 0 {
		fmt.Println("одинаковые файлы не обнаружены")
		return
	}

	if !removeCopies {
		return // ключ удаления копий не задан
	}

	// спросить пользователя, хочет он удалять копии или нет
	var answer string
	fmt.Print("удалить копии (Y/N)? ")
	fmt.Scanln(&answer)
	if answer != "y" && answer != "Y" {
		return
	}

	for _, paths := range copies {
		// получать пути к файлам начиная со второго!!!
		for i := 1; i < len(paths); i++ {
			err := os.Remove(paths[i])
			if err != nil {
				log.Println(err)
			} else {
				fmt.Println("удален файл: ", paths[i])
			}
		}
	}
}

// IterateEntitiesInDirectory перебирает файлы в каталоге
// и его подкаталогах начиная со startPath
// Полученные пути к файлам отправляет в канал filePathChan
func IterateEntitiesInDirectory(startPath string, filePathChan chan string) {
	entries, err := os.ReadDir(startPath)
	if err != nil {
		log.Fatal(err)
	}

	for _, entry := range entries {
		curPath := fmt.Sprintf("%s/%s", startPath, entry.Name())
		if entry.IsDir() {
			// вызвать рекурсивно для подкаталога
			IterateEntitiesInDirectory(curPath, filePathChan)
		} else {
			// отправить путь к файлу в канал
			filePathChan <- curPath
		}

	}
}
