package main

import (
	"github.com/sirupsen/logrus"
	"lesson3/pkg/config"

	f "lesson3/pkg/files"
	p "lesson3/pkg/program"
)

func main() {
	log := logrus.New()
	customFormatter := new(logrus.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetLevel(logrus.DebugLevel)
	log.SetFormatter(customFormatter)

	log.Info("Service started")
	defer log.Info("Service finished")

	log.Info("Config initialization started")
	cnfg, err := config.NewAppConfig()
	if err != nil {
		log.WithError(err).Warn("Invalid config set. Process was stopped")
		return
	}
	log.Info("Config initialization completed")

	uniqueFiles := f.NewUniqueFilesMap(log)

	program := p.NewProgram(
		cnfg,
		uniqueFiles,
		uniqueFiles,
		uniqueFiles,
		uniqueFiles,
		log,
	)

	err = program.Start()
	if err != nil {
		log.WithError(err).Fatal("Failed to process")
	}
}
