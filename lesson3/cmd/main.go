package main

import (
	"lesson3/cmd/internal/app"
	"lesson3/cmd/internal/pkg/config"
	"lesson3/cmd/internal/pkg/files"
	"lesson3/cmd/internal/pkg/log"
)

func main() {
	log := log.NewLogWithConfuguration()
	log.Info("Service started")
	defer log.Info("Service finished")

	log.Info("Config initialization started")
	cnfg, err := config.NewAppConfig()
	if err != nil {
		log.WithError(err).Warn("Invalid config set. Process was stopped")
		return
	}
	log.Info("Config initialization completed")

	uniqueFiles := files.NewUniqueFilesMap(log)

	program := app.NewService(
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
