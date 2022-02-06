//Добавьте в программу для поиска дубликатов (либо же краулер разработанный на первом занятии),
//разработанную в рамках проектной работы на предыдущем модуле логи.
//Необходимо использовать пакет zap или logrus.
//Разграничить уровни логирования.
//Обогатить параметрами по вашему усмотрению.
//Вставить вызов panic() в участке коде, в котором осуществляется переход в поддиректорию;
//удостовериться, что по логам можно локализовать при каком именно переходе в какую директорию сработала паника
//Для краулера что бы можно было понять на каком url произошло

package main

import (
	"github.com/sirupsen/logrus"

	"lesson2/pkg/config"
	f "lesson2/pkg/files"
	p "lesson2/pkg/program"
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

	program := p.NewProgram(cnfg, uniqueFiles, log)
	err = program.Start()
	if err != nil {
		log.WithError(err).Fatal("Failed to process")
	}
}
