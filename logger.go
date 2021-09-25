package main

import (
	"blog/config"
	"log"

	"go.uber.org/zap/zapcore"

	"go.uber.org/zap"
)

//getLogger creates a new instance of the logger based on the config
func getLogger(production bool, cfg *config.Logger) *zap.Logger {
	if production {
		logger, err := zap.NewProduction()
		if err != nil {
			log.Panicf("Failed to initialize production logger: %v", err)
		}

		return logger
	}

	devConfig := zap.NewDevelopmentConfig()
	devConfig.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder

	logger, err := devConfig.Build()
	if err != nil {
		log.Panicf("Failed to initialize development logger: %v", err)
	}

	return logger
}
