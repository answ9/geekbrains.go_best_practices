//Package config provides AppConfig struct that reads parameters from the command line arguments (flags) and validates it
package config

import (
	"flag"
	"fmt"
)

//Flags that can be used by user for app configuration
var (
	path             = flag.String("path", ".", "a path for the app to find dublicates of files")
	workers          = flag.Int("workers", 5, "amount of workers")
	deleteDublicates = flag.Bool("delete", false, "delete the found dublicates?")
	printResult      = flag.Bool("print-result", true, "print the list of found files and duplicates in console?")
)

//AppConfig contains configuration parameters defined by user or set by default
type AppConfig struct {
	Path             string
	Workers          int
	DeleteDublicates bool
	PrintResult      bool
}

//Method Validate() validates set configuration and returns an error if it fails
func (c *AppConfig) Validate() error {
	if c.Workers < 1 || c.Workers > 50 {
		return fmt.Errorf("Amount of workers is limited from 1 to 50")
	}
	if c.Path == "" {
		return fmt.Errorf("Path cant be empty")
	}

	return nil
}

// Use method NewAppConfig() to create a new AppConfig
func NewAppConfig() (*AppConfig, error) {
	flag.Parse()
	config := &AppConfig{*path, *workers, *deleteDublicates, *printResult}
	return config, config.Validate()
}
