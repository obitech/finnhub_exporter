package main

import (
	"go.uber.org/zap"
)

func newLogger(level string) (*zap.Logger, error) {
	atom := zap.NewAtomicLevel()

	switch level {
	case "debug":
		atom.SetLevel(zap.DebugLevel)
	case "warn":
		atom.SetLevel(zap.WarnLevel)
	case "error":
		atom.SetLevel(zap.ErrorLevel)
	default:
		level = "info"
		atom.SetLevel(zap.InfoLevel)
	}

	cfg := zap.Config{
		Development:       false,
		DisableCaller:     true,
		DisableStacktrace: false,
		EncoderConfig:     zap.NewProductionEncoderConfig(),
		Encoding:          "json",
		ErrorOutputPaths:  []string{"s"},
		Level:             atom,
		OutputPaths:       []string{"stdout"},
	}

	return cfg.Build()
}
