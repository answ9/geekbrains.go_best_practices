package main

import (
	a "task_2_3_4/internal/app"
	c "task_2_3_4/internal/pkg/config"
	f "task_2_3_4/internal/pkg/files"
	l "task_2_3_4/internal/pkg/log"
)

func main() {
	log := l.NewLogWithConfuguration()
	log.Info("Service started")
	defer log.Info("Service finished")

	log.Info("Config initialization started")
	cnfg, err := c.NewAppConfig()
	if err != nil {
		log.WithError(err).Warn("Invalid config set. Process was stopped")
		return
	}
	log.Info("Config initialization completed")

	uniqueFiles := f.NewUniqueFilesMap(log)

	program := a.NewService(
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
