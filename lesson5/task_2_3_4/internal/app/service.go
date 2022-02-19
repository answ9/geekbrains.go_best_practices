// Package provides the main program/app work logic hidden behind a "facade"
// linter: commentFormatting: put a space between `//` and comment text (gocritic)

package app

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type ConfigGetter interface {
	Get() (string, int, bool, bool)
}

type Finder interface {
	Find(string, int) int
}

type Printer interface {
	Print(string)
}

type Deleter interface {
	DeleteDuplicates() error
}

type CountGetter interface {
	GetDuplicatesCount() int
}

// Struct Service consists of Config, UniqueFiles and amount of found file duplicates
type Service struct {
	Config  ConfigGetter
	Finder  Finder
	Printer Printer
	Deleter Deleter
	Getter  CountGetter
	log     *logrus.Logger
}

// Method Start() is used to start the program finding all the files and their duplicates in a given directory.
// The program prints in console the list of found files and their duplicates.
// It additionally asks the user to confirm deletion if the parameter "DeleteDublicates" is in true.
func (p Service) Start() error {
	const yes, no = "yes", "no"
	searchPath, workers, printResult, deleteDublicates := p.Config.Get()

	duplicates := p.Finder.Find(searchPath, workers)
	if printResult {
		p.Printer.Print(searchPath)
	}

	if !deleteDublicates || duplicates == 0 {
		return nil
	}
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Are you sure to delete all duplicate files? yes/no: ")
	for scanner.Scan() {
		if scanner.Err() != nil {
			return scanner.Err()
		}
		in := strings.TrimSpace(scanner.Text()) //linter: string `yes` has 2 occurrences, make it a constant (goconst)
		if in != yes && in != no {
			fmt.Printf("%s", "try again: type yes or no\n")
			continue
		}
		if in != yes {
			p.log.Info("User did not confirm deletion")
			break
		}

		p.log.Info("User confirmed deletion")
		err := p.Deleter.DeleteDuplicates()
		if err != nil {
			p.log.WithError(err).Error("Failed to delete the duplicates")
			return err
		}
		p.log.Info("All duplicate files were deleted")
		break
	}

	return nil
}

// Use method NewService() to get a new program to start
func NewService(c ConfigGetter, f Finder, p Printer, d Deleter, g CountGetter, l *logrus.Logger) Service {
	return Service{c, f, p, d, g, l}
}
