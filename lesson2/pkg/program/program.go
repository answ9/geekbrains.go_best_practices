//Package provides the main program/app work logic hidden behind a "facade"
package program

import (
	"bufio"
	"crypto/sha512"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/TwiN/go-color"
	"github.com/sirupsen/logrus"

	"lesson2/pkg/config"
	f "lesson2/pkg/files"
)

// Struct Program consists of Config, UniqueFiles and amount of found file duplicates
type Program struct {
	Config      *config.AppConfig
	UniqueFiles *f.UniqueFiles
	Duplicates  int
	log         *logrus.Logger
}

// Method Start() is used to start the program finding all the files and their duplicates in a given directory.
// The program prints in console the list of found files and their duplicates.
// It additionally asks the user to confirm deletion if the parameter "DeleteDublicates" is in true.
func (p *Program) Start() error {
	files := make(chan f.File)

	go func(dir string, files chan<- f.File) {
		filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				p.log.ReportCaller = true
				p.log.WithError(err).Fatal(fmt.Sprintf("Error during filepath walk %s", fmt.Sprintf("%s/%s", dir, path)))
			}
			if !info.IsDir() && info.Name() != ".DS_Store" {
				files <- f.NewFile(path, info.Name())
			}
			return nil
		})
		close(files)
	}(p.Config.Path, files)

	var wg sync.WaitGroup
	wg.Add(p.Config.Workers)

	for i := 0; i < p.Config.Workers; i++ {
		func(files <-chan f.File, uniqueFiles *f.UniqueFiles, wg *sync.WaitGroup) {
			for file := range files {
				data, err := ioutil.ReadFile(path.Join(".", file.Path))
				if err != nil {
					p.log.WithError(err).Error("failed to read the file", file.Path)
				}
				digest := sha512.Sum512(data)
				uniqueFiles.Mtx.Lock()
				if _, ok := uniqueFiles.Map[digest]; ok {
					p.Duplicates++
				}
				uniqueFiles.Map[digest] = append(uniqueFiles.Map[digest], file)
				uniqueFiles.Mtx.Unlock()
			}
			wg.Done()
		}(files, p.UniqueFiles, &wg)
	}
	wg.Wait()

	p.UniqueFiles.Sort()
	p.printResult()

	err := p.askForConfirmBeforeDeletion()
	if err != nil {
		p.log.WithError(err).Error("failed to confirm the deletion of duplicates")
		return err
	}

	return nil
}

// Method printResult() is used inside Start() and prints in console the list of found files and their duplicates
func (p *Program) printResult() {
	if !p.Config.PrintResult {
		return
	}
	p.log.Info(fmt.Sprintf("Found %d unique files and %d duplicates in \"%s\":\n", len(p.UniqueFiles.Map), p.Duplicates, p.Config.Path))

	counter := 1
	for k, _ := range p.UniqueFiles.Map {
		txt := ""
		for i, _ := range p.UniqueFiles.Map[k] {
			if i == 0 {
				txt += fmt.Sprintf("%d) file %s", counter, color.Ize(color.Blue, p.UniqueFiles.Map[k][i].Name))
				if len(p.UniqueFiles.Map[k]) > 1 {
					txt += fmt.Sprintf(" with %d duplicates:", len(p.UniqueFiles.Map[k])-1)
				}

			} else {
				txt += fmt.Sprintf(" %s |", color.Ize(color.Blue, p.UniqueFiles.Map[k][i].Name))
			}
		}
		p.log.Info(txt)
		counter++
	}
}

// Method askForConfirmBeforeDeletion() is used inside Start() and contains logic of deleting duplicates files if such were found and user wanted to delete them
func (p *Program) askForConfirmBeforeDeletion() error {
	if p.Config.DeleteDublicates && p.Duplicates > 0 {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Print("Are you sure to delete all duplicate files? yes/no: ")

		for scanner.Scan() {
			if scanner.Err() != nil {
				return scanner.Err()
			}
			in := strings.TrimSpace(scanner.Text())
			if in != "yes" && in != "no" {
				fmt.Printf("%s", "try again: type yes or no\n")
				continue
			}
			if in != "yes" {
				p.log.Info("User did not confirm deletion")
				break
			}

			p.log.Info("User confirmed deletion")
			err := p.UniqueFiles.DeleteDuplicates()
			if err != nil {
				p.log.WithError(err).Error("Failed to delete the duplicates")
				return err
			}
			p.log.Info("All duplicate files were deleted")
			break
		}
	}

	return nil
}

// Use method NewProgram() to get a new program to start
func NewProgram(cnfg *config.AppConfig, uniqueFiles *f.UniqueFiles, log *logrus.Logger) *Program {
	return &Program{cnfg, uniqueFiles, 0, log}
}
