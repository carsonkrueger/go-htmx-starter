package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/carsonkrueger/main/cfg"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(cfg *cfg.Config) *zap.Logger {
	var config zap.Config

	switch strings.ToLower(cfg.AppEnv) {
	case "development":
		config = zap.Config{
			Level:         zap.NewAtomicLevelAt(zap.DebugLevel),
			Development:   true,
			Encoding:      "json",
			EncoderConfig: zap.NewDevelopmentEncoderConfig(),
			OutputPaths: []string{
				"stdout",
			},
		}
	case "production":
		config = zap.Config{
			Level:         zap.NewAtomicLevelAt(zap.ErrorLevel),
			Development:   false,
			Encoding:      "json",
			EncoderConfig: zap.NewProductionEncoderConfig(),
			OutputPaths: []string{
				"stderr",
			},
			InitialFields: map[string]any{
				"pid": os.Getpid(),
			},
		}
	default:
		panic(fmt.Sprintf("Invalid app environment: %s", cfg.AppEnv))
	}

	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	return zap.Must(config.Build())
}
