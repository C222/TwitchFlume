package logger

import "github.com/Sirupsen/logrus"

var log *logrus.Logger = nil

func Init()  {
	log = logrus.New()
	log.Formatter = new(logrus.TextFormatter)
}

func GetLogger() *logrus.Logger {
	if log == nil {
		Init()
	}
	return log
}