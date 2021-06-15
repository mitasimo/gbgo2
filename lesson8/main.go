package main

import (
	"crypto/md5"
	"flag"
	"fmt"
	"hash"
	"io"
	"log"
	"os"
	"runtime"
	"sync"
)

// FileHash связывает путь к файлу и его хэш
type FileHash struct {
	hash.Hash
	FilePath string
}

func main() {

	var (
		startPath  string // каталог с которого начинать перебор файлов
		needRemove bool   // признак необходимости удаления копий
	)

	flag.StringVar(&startPath, "p", ".", "start path")
	flag.BoolVar(&needRemove, "r", false, "remove copies")

	flag.Parse()

	if startPath == "" {
		log.Fatal("start path is not set")
	}

	var (
		wg sync.WaitGroup
	)

	filePathChan := make(chan string)
	fileHashChan := make(chan *FileHash)

	// запустить горутину, считывающую имена файлав из каталога
	go func() {
		entries, err := os.ReadDir(startPath)
		if err != nil {
			log.Fatal(err)
		}

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			// добавлять пути к файлам в канал
			filePathChan <- startPath + "/" + entry.Name()
		}

		// как файлы закончислись, закрыть канал
		close(filePathChan)
		// close(doneChan)
	}()

	// запустить горутины подсчета хеша файлов
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for filePath := range filePathChan {
				fileHash := &FileHash{
					FilePath: filePath,
					Hash:     md5.New(),
				}

				file, err := os.Open(filePath)
				if err != nil {
					log.Println(err)
					continue
				}

				io.Copy(fileHash, file)
				fileHashChan <- fileHash
			}
		}()
	}

	go func() {
		wg.Wait()
		close(fileHashChan)
	}()

	// мапа для сбора данных о дублях
	// ключ - хеш в виде строки
	// значение - массив путей к файлам
	doubles := make(map[string][]string)

	// читаем из канала хеши
	for fileHash := range fileHashChan {
		// добавляем к хешам пути
		strHash := fmt.Sprintf("%x", fileHash.Sum(nil))
		filesPath := doubles[strHash]
		filesPath = append(filesPath, fileHash.FilePath)
		doubles[strHash] = filesPath
	}

	if len(doubles) == 0 {
		fmt.Println("doubles not found")
		return
	}

	doublesPrinted := 0

	for key, pathes := range doubles {
		if len(pathes) > 1 {
			doublesPrinted++
			fmt.Println(key)
			for _, curPath := range pathes {
				fmt.Println("\t", curPath)
			}
		}
	}

	if doublesPrinted == 0 {
		fmt.Println("doubles not found")
	}

}
