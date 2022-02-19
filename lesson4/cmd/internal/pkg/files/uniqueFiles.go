package files

import (
	"crypto/sha512"
	"fmt"
	"github.com/TwiN/go-color"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"sync"
)

// Struct UniqueFiles contains a protected by Mutex map with slices of duplicated files.
// Duplicates must be determined according to hash comparison
type UniqueFiles struct {
	Mtx        *sync.Mutex
	Map        map[[sha512.Size]byte][]File
	log        *logrus.Logger
	duplicates int
}

// Method GetDuplicatesCount() returns the count of duplicates
func (uf *UniqueFiles) GetDuplicatesCount() int {
	return uf.duplicates
}

// Method Find() walks through the path and finds all files and their duplicates
func (uf *UniqueFiles) Find(searchPath string, workers int) int {
	files := make(chan File)

	go func(dir string, files chan<- File) {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				uf.log.ReportCaller = true
				uf.log.WithError(err).Fatal(fmt.Sprintf("Error during filepath walk %s", fmt.Sprintf("%s/%s", dir, path)))
			}
			if !info.IsDir() && info.Name() != ".DS_Store" {
				files <- NewFile(path, info.Name())
			}
			return nil
		})
		close(files)
	}(searchPath, files)

	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		func(files <-chan File, uniqueFiles *UniqueFiles, wg *sync.WaitGroup) {
			for file := range files {
				data, err := ioutil.ReadFile(path.Join(".", file.Path))
				if err != nil {
					uf.log.WithError(err).Error("failed to read the file", file.Path)
				}
				digest := sha512.Sum512(data)
				uniqueFiles.Mtx.Lock()
				if _, ok := uniqueFiles.Map[digest]; ok {
					uf.duplicates++
				}
				uniqueFiles.Map[digest] = append(uniqueFiles.Map[digest], file)
				uniqueFiles.Mtx.Unlock()
			}
			wg.Done()
		}(files, uf, &wg)
	}
	wg.Wait()
	uf.Sort()
	return uf.duplicates
}

// Method Print() prints list of files and their duplicates
func (uf *UniqueFiles) Print(searchPath string) {
	uf.log.Info(fmt.Sprintf("Found %d unique files and %d duplicates in \"%s\":\n", len(uf.Map), uf.duplicates, searchPath))
	counter := 1
	for k, _ := range uf.Map {
		txt := ""
		for i, _ := range uf.Map[k] {
			if i == 0 {
				txt += fmt.Sprintf("%d) file %s", counter, color.Ize(color.Blue, uf.Map[k][i].Name))
				if len(uf.Map[k]) > 1 {
					txt += fmt.Sprintf(" with %d duplicates:", len(uf.Map[k])-1)
				}

			} else {
				txt += fmt.Sprintf(" %s |", color.Ize(color.Blue, uf.Map[k][i].Name))
			}
		}
		uf.log.Info(txt)
		counter++
	}
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
	return &UniqueFiles{&sync.Mutex{}, make(map[[sha512.Size]byte][]File), log, 0}
}
