package fs

import (
	"os"
	"strings"
)

// FileSystemEntry описывает функционал сущности файловой системы
type FileSystemEntry interface {
	Path() string
	IsDir() bool
	SubEntries() ([]FileSystemEntry, error)
}

func New(startPath string) *FileSystemEntryImpl {

}

// FileSystemEntryImpl релизует функционал...
type FileSystemEntryImpl struct {
	path  string
	isDir bool
}

func (fsi *FileSystemEntryImpl) Path() string {
	return fsi.path
}
func (fsi *FileSystemEntryImpl) IsDir() bool {
	return fsi.isDir
}

func (fsi *FileSystemEntryImpl) SubEntries() ([]FileSystemEntry, error) {
	subEntries := make([]FileSystemEntry, 0)

	entries, err := os.ReadDir(fsi.path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		subEntries = append(subEntries, &FileSystemEntryImpl{
			path:  joinStrings(fsi.path, "/", entry.Name()),
			isDir: entry.IsDir(),
		})
	}

	return subEntries, nil
}

func joinStrings(strs ...string) string {
	return strings.Join(strs, "/")
}
