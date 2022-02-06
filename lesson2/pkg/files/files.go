// Package files is represented by two structs: File and UniqueFiles.
package files

import (
	"crypto/sha512"
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/sirupsen/logrus"
)

// Struct File contains Path and Name of a file.
// Path can be used for fast access of a definite file
// Name can be used mostly for printing it to user
type File struct {
	Path, Name string
}

// Use method NewFile(path, name string) to create a new object
func NewFile(path, name string) File {
	return File{path, name}
}

// Struct UniqueFiles contains a protected by Mutex map with slices of duplicated files.
// Duplicates must be determined according to hash comparison
type UniqueFiles struct {
	Mtx *sync.Mutex
	Map map[[sha512.Size]byte][]File
	log *logrus.Logger
}

// Method Sort() sorts file duplicates in slice according to length of names.
// Thus, the first one in the slice with the shortest name is considered to be the original one, and all the others are considered to be duplicates.
func (uf *UniqueFiles) Sort() {
	for k, _ := range uf.Map {
		if len(uf.Map[k]) == 1 {
			continue
		}
		sort.Slice(uf.Map[k], func(i, j int) bool { return len(uf.Map[k][i].Name) < len(uf.Map[k][j].Name) })
	}
}

// Method DeleteDuplicates() loops over slice of duplicates files and deletes those that considered to be duplicates (all except the first in a slice).
// Is recommended to use after method Sort().
func (uf *UniqueFiles) DeleteDuplicates() error {
	uf.log.Info("Duplicate files deletion started")
	for k, _ := range uf.Map {
		if len(uf.Map[k]) == 1 {
			continue
		}
		for i, _ := range uf.Map[k] {
			if i == 0 {
				continue
			}
			uf.log.Info(fmt.Sprintf("...deleting %s", uf.Map[k][i].Path))
			err := os.Remove(uf.Map[k][i].Path)
			if err != nil {
				uf.log.WithError(err).Warn("Failed to delete file", uf.Map[k][i].Path)
				return err
			}
		}
	}
	uf.log.Info("Duplicate files deletion ended")
	return nil
}

// Use method NewUniqueFilesMap() to create a protected by Mutex map with slices of duplicated files
func NewUniqueFilesMap(log *logrus.Logger) *UniqueFiles {
	return &UniqueFiles{&sync.Mutex{}, make(map[[sha512.Size]byte][]File), log}
}
