package config

import (
	"flag"
	"fmt"
	"net/url"
	"strings"
)

var (
	startUrl  = flag.String("startUrl", "https://www.w3.org/Consortium/", "url page for the crawler to start with")
	curDepth  = flag.Int("curDepth", 0, "current depth")
	maxDepth  = flag.Int("maxDepth", 1, "max depth")
	maxErrors = flag.Int("maxErrors", 5, "max errors count before quit")
	timeOut   = flag.Int64("timeOut", 10, "context time out")
)

type AppConfig struct {
	StartUrl  string
	CurDepth  int
	MaxDepth  int
	MaxErrors int
	TimeOut   int64
}

func (c *AppConfig) Validate() error {
	if strings.TrimSpace(c.StartUrl) == "" {
		return fmt.Errorf("Your start url cant be empty")
	}

	_, err := url.ParseRequestURI(c.StartUrl)
	if err != nil {
		return fmt.Errorf("Start url is invalid %v", err)
	}

	if c.CurDepth < 0 || c.CurDepth >= c.MaxDepth {
		return fmt.Errorf("Current depth cant be negative and equal or bigger than max depth")
	}
	if c.MaxDepth < 1 || c.MaxDepth > 10 {
		return fmt.Errorf("Max depth is limited from 1 to 10")
	}

	if c.TimeOut < 1 {
		return fmt.Errorf("Time out cant be equal or smaller than 0")
	}

	if c.MaxErrors < 1 || c.MaxErrors > 99 {
		return fmt.Errorf("Errors are limited from 1 to 99")
	}

	return nil
}

func NewAppConfig() (*AppConfig, error) {
	flag.Parse()
	config := &AppConfig{*startUrl, *curDepth, *maxDepth, *maxErrors, *timeOut}
	return config, config.Validate()
}
