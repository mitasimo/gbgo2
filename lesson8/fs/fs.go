package fs

import (
	"errors"
	"io"
	"os"
	"strings"
)

var (
	ErrOutOfBound = errors.New("out of bound")
)

func New(startPath string, recur bool) (*FileSystemIterator, error) {
	paths, err := iterateFilesInDirerctory(startPath, recur)
	if err != nil {
		return nil, err
	}

	return &FileSystemIterator{
		currentIndex: -1,
		paths:        paths,
	}, nil
}

// FileSystemEntryImpl релизует функционал...
type FileSystemIterator struct {
	currentIndex int
	paths        []string
}

func (fsi *FileSystemIterator) Next() bool {
	if len(fsi.paths) == 0 || fsi.currentIndex >= len(fsi.paths)-1 {
		return false
	}
	fsi.currentIndex++
	return true
}

func (fsi *FileSystemIterator) Path() (string, error) {
	if fsi.currentIndex < 0 {
		return "", ErrOutOfBound
	}
	return fsi.paths[fsi.currentIndex], nil
}

func (fsi *FileSystemIterator) ReadCloser() (io.ReadCloser, error) {
	if fsi.currentIndex < 0 {
		return nil, ErrOutOfBound
	}
	return os.Open(fsi.paths[fsi.currentIndex])
}

func iterateFilesInDirerctory(startPath string, recur bool) ([]string, error) {
	paths := make([]string, 0)

	entries, err := os.ReadDir(startPath)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {

		newPath := strings.Join(StringArgsArrTo(startPath, entry.Name()), "/")

		if entry.IsDir() {
			if !recur {
				continue
			}

			subPath, err := iterateFilesInDirerctory(newPath, recur)
			if err != nil {
				return nil, err
			}
			paths = append(paths, subPath...)
		} else {
			paths = append(paths, newPath)
		}
	}

	return paths, nil
}

func StringArgsArrTo(strs ...string) []string {
	return strs
}
