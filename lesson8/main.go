// Утилита для поиска и удаления задублированных файлов
// в каталоге и его подкаталогах
// Вызывается
// $ util -p=<стартовый каталог> -r
// ключ -dir (directory) задает стартовый каталог (по уполчанию текущий каталог программы)
// ключ -rm (remove) заставит утилиту удалить файлы (оставит только один)

package main

import (
	"errors"
	"flag"
	"fmt"
	"hash/adler32"
	"io"
	"os"
	"runtime"
	"sync"

	log "github.com/sirupsen/logrus"
)

// FileHash связывает путь к файлу и его хэш
type FileHash struct {
	Hash     uint32
	FilePath string
	Err      error
}

func main() {

	logger := log.New()
	logger.SetFormatter(&log.JSONFormatter{})

	var (
		startPath    string // каталог с которого начинать перебор файлов
		removeCopies bool   // признак необходимости удаления копий
	)

	// чтение флогов командной строки
	flag.StringVar(&startPath, "dir", ".", "начальный каталог")
	flag.BoolVar(&removeCopies, "rm", false, "удалять копии")
	flag.Parse()

	if startPath == "" {
		logger.Error("не задан начальный каталог")
	}

	logger.WithField("dir", startPath).Info("обрабатывается каталог")

	var (
		wg sync.WaitGroup
	)

	filePathChan := make(chan string)    // канал для передачи путей к файлам
	fileHashChan := make(chan *FileHash) // канал для передачи хеша файлов

	// запуск горутину=ы, считывающей имена файлав из каталога
	go func() {
		// после завершения итераций по файлам закрыть канал путей к файлам
		defer close(filePathChan)
		IterateEntitiesInDirectory(startPath, filePathChan, logger)
	}()

	// запуск горутин подсчета хеша файлов
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for filePath := range filePathChan { // получить очередной путь к файлу из канала
				fileHash := &FileHash{FilePath: filePath}
				// открыть файл по пути
				file, err := os.Open(filePath)
				if err != nil {
					//log.Println(err) // записать ошибку в лог
					fileHash.Err = err
				} else {
					// хешировать данные файла
					hash, err := CalculateAdler32Hash(file)
					file.Close()

					if err != nil {
						fileHash.Err = err
					} else {
						fileHash.Hash = hash
					}
				}
				// отправить хеш файла в канал
				fileHashChan <- fileHash
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
		ll := logger.WithField("Path", fileHash.FilePath)
		if fileHash.Err != nil {
			ll.Errorf("ошибка получения хеша: %v", fileHash.Err)
		} else {
			filesPath := copies[fileHash.Hash]               // получить массив путей к файлам
			filesPath = append(filesPath, fileHash.FilePath) // добавить путь к массиву
			copies[fileHash.Hash] = filesPath                // сохранить новый массив путей

			ll.Info("хеш посчитан")
		}
	}

	copiesPrinted := 0 // количество напечатанных путей к дублям

	for key, pathes := range copies {
		if len(pathes) > 1 { // есть копии
			copiesPrinted++
			ll := logger.WithField("hash", key)
			//fmt.Println("Хеш:", key) // вывести хеш
			for _, curPath := range pathes {
				ll.WithField("path", curPath).Info("equal") // вывести путь к файлу
			}
		}
	}

	if copiesPrinted == 0 {
		logger.Info("одинаковые файлы не обнаружены")
		return
	}

	if !removeCopies {
		return // ключ удаления копий не задан
	}

	// спросить пользователя, хочет он удалять копии или нет
	// var answer string
	// fmt.Print("удалить копии (Y/N)? ")
	// fmt.Scanln(&answer)
	// if answer != "y" && answer != "Y" {
	// 	return
	// }

	for _, paths := range copies {
		// получать пути к файлам начиная со второго!!!
		for i := 1; i < len(paths); i++ {
			ll := logger.WithField("Path", paths[i])
			err := os.Remove(paths[i])
			if err != nil {
				ll.Errorf("ошибка удаления файла: %v", err)
			} else {
				ll.Info("файл удален")
			}
		}
	}
}

// IterateEntitiesInDirectory перебирает файлы в каталоге
// и его подкаталогах начиная со startPath
// Полученные пути к файлам отправляет в канал filePathChan
func IterateEntitiesInDirectory(startPath string, filePathChan chan string, logger *log.Logger) {
	ll := logger.WithField("dir", startPath)

	entries, err := os.ReadDir(startPath)
	if err != nil {
		ll.Errorf("ошибка чтения каталога")
		return
		//log.Fatal(err)
	}
	ll.Info("обрабатывается каталог")

	for _, entry := range entries {
		curPath := fmt.Sprintf("%s/%s", startPath, entry.Name())
		if entry.IsDir() {
			// вызвать рекурсивно для подкаталога
			IterateEntitiesInDirectory(curPath, filePathChan, logger)
		} else {
			// отправить путь к файлу в канал
			filePathChan <- curPath
		}
	}
}

func CalculateAdler32Hash(r io.Reader) (uint32, error) {
	if r == nil {
		return 0, errors.New("nil reader")
	}
	hash := adler32.New()
	_, err := io.Copy(hash, r)
	if err != nil && err != io.EOF {
		return 0, err
	}
	return hash.Sum32(), nil
}
