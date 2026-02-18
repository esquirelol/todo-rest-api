package logger

import (
	"log"

	"go.uber.org/zap"
)

func LoggerInit(env string) *zap.Logger {
	var (
		logger *zap.Logger
		err    error
	)
	switch env {
	case "development":
		logger, err = zap.NewDevelopment()
		if err != nil {
			log.Fatal(err.Error())
		}
	case "production":
		logger, err = zap.NewProduction()
		if err != nil {
			log.Fatal(err.Error())
		}

	default:
		logger = zap.NewExample()
	}
	return logger

}
