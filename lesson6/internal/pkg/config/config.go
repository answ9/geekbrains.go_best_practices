// Package config provides AppConfig struct that reads parameters from the command line arguments (flags) and validates it
// linter: commentFormatting: put a space between `//` and comment text (gocritic)
package config

import (
	"flag"
	"fmt"
)

// Flags that can be used by user for app configuration
// linter: commentFormatting: put a space between `//` and comment text (gocritic)
var (
	path             = flag.String("path", "../..", "a path for the app to find dublicates of files")
	workers          = flag.Int("workers", 5, "amount of workers")
	deleteDublicates = flag.Bool("delete", false, "delete the found dublicates?")
	printResult      = flag.Bool("print-result", true, "print the list of found files and duplicates in console?")
)

// AppConfig contains configuration parameters defined by user or set by default
// linter: commentFormatting: put a space between `//` and comment text (gocritic)
type AppConfig struct {
	Path             string
	Workers          int
	DeleteDublicates bool
	PrintResult      bool
}

// Method Validate() validates set configuration and returns an error if it fails
// linter: commentFormatting: put a space between `//` and comment text (gocritic)
func (c *AppConfig) Validate() error {
	if c.Workers < 1 || c.Workers > 50 {
		return fmt.Errorf("amount of workers is limited from 1 to 50") // linter: ST1005: error strings should not be capitalized (stylecheck)
		// linter: unnecessary trailing newline (whitespace)
	}
	if c.Path == "" {
		return fmt.Errorf("path cant be empty") // linter: ST1005: error strings should not be capitalized (stylecheck)
	}

	return nil
}

func (c *AppConfig) Get() (string, int, bool, bool) {
	return c.Path, c.Workers, c.PrintResult, c.DeleteDublicates
}

// Use method NewAppConfig() to create a new AppConfig
func NewAppConfig() (*AppConfig, error) {
	flag.Parse()
	config := &AppConfig{*path, *workers, *deleteDublicates, *printResult}
	return config, config.Validate()
}
