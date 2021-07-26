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
	"hash/adler32"
	"io"
	"os"
	"runtime"
	"sync"

	"github.com/mitasimo/gbgo2/lesson8/fs"
	log "github.com/sirupsen/logrus"
)

// FilesIterator описывает функционал итератор файлов
type FilesIterator interface {
	// Next - переходит к следующему файл
	Next() bool
	// Path - возвращает путь к текущему файлу
	Path() (string, error)
	// Reader - возвращает io.ReadCloser текущего файла или ошибку
	ReadCloser() (io.ReadCloser, error)
}

// FileHash связывает путь к файлу и его хэш
type FileHash struct {
	Hash uint32
	Path string
	Err  error
}

type FileEntry struct {
	Path string
	RC   io.ReadCloser
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

	fi, err := fs.New(startPath, true)
	if err != nil {
		logger.Error(err)
		os.Exit(-1)
	}

	logger.WithField("dir", startPath).Debug("обрабатывается каталог")

	var (
		wg sync.WaitGroup
	)

	fileEntryChan := make(chan *FileEntry) // канал для передачи путей к файлам
	fileHashChan := make(chan *FileHash)   // канал для передачи хеша файлов

	// запуск горутины, считывающей имена файлав из каталога
	go func() {
		// после завершения итераций по файлам закрыть канал путей к файлам
		defer close(fileEntryChan)
		IterateEntitiesInDirectory(fi, fileEntryChan, logger)
	}()

	// запуск горутин подсчета хеша файлов
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()

			for fileEntry := range fileEntryChan { // получить очередной путь к файлу из канала
				fileHash := &FileHash{Path: fileEntry.Path}

				// хешировать данные файла
				hash, err := CalculateAdler32Hash(fileEntry.RC)
				fileEntry.RC.Close()
				if err != nil {
					fileHash.Err = err
				} else {
					fileHash.Hash = hash
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
		ll := logger.WithField("Path", fileHash.Path)
		if fileHash.Err != nil {
			ll.Errorf("ошибка получения хеша: %v", fileHash.Err)
		} else {
			filesPath := copies[fileHash.Hash]           // получить массив путей к файлам
			filesPath = append(filesPath, fileHash.Path) // добавить путь к массиву
			copies[fileHash.Hash] = filesPath            // сохранить новый массив путей

			ll.Debug("хеш посчитан")
		}
	}

	copiesPrinted := 0 // количество напечатанных путей к дублям

	for key, pathes := range copies {
		if len(pathes) > 1 { // есть копии
			copiesPrinted++
			ll := logger.WithField("hash", key)
			ll.Info("duplicated hash")
			//fmt.Println("Хеш:", key) // вывести хеш
			for _, curPath := range pathes {
				ll.WithField("path", curPath).Info("duplicated file") // вывести путь к файлу
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
func IterateEntitiesInDirectory(fi FilesIterator, fileEntryChan chan *FileEntry, logger *log.Logger) {
	for fi.Next() {
		path, _ := fi.Path()
		rc, err := fi.ReadCloser()
		if err != nil {
			logger.WithField("Path", path).Error(err.Error(), "get GetReadCloser")
		} else {
			fileEntryChan <- &FileEntry{
				Path: path,
				RC:   rc,
			}
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
