package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// CreateLogger initializes a zap.Logger based on the provided environment.
// It configures the logger for production or development environments,
// setting appropriate encoding, output paths, and encoder configurations.
// The logger includes a timestamp and process ID in its initial fields.
func CreateLogger(env string) *zap.Logger {
	var config zap.Config
	var encoderCfg zapcore.EncoderConfig

	if env == "prod" || env == "production" {
		encoderCfg = zap.NewProductionEncoderConfig()
	} else {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	}

	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	switch env {
	case "prod", "production":
		config = zap.Config{
			Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
			Development:       false,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "json",
			EncoderConfig:     encoderCfg,
			OutputPaths:       []string{"stderr"},
			ErrorOutputPaths:  []string{"stderr"},
			InitialFields:     map[string]interface{}{"pid": os.Getpid()},
		}
	default:
		config = zap.Config{
			Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
			Development:       true,
			DisableCaller:     false,
			DisableStacktrace: false,
			Encoding:          "console",
			EncoderConfig:     encoderCfg,
			OutputPaths:       []string{"stdout"},
			ErrorOutputPaths:  []string{"stderr"},
			InitialFields:     map[string]interface{}{"pid": os.Getpid()},
		}
	}

	return zap.Must(config.Build())
}
