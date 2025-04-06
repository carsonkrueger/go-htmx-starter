package logger

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/carsonkrueger/main/cfg"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func NewLogger(cfg *cfg.Config) *zap.Logger {
	var config zap.Config

	switch strings.ToLower(cfg.AppEnv) {
	case "development":
		config = zap.Config{
			Level:             zap.NewAtomicLevelAt(zap.DebugLevel),
			Development:       true,
			Encoding:          "console",
			EncoderConfig:     zap.NewDevelopmentEncoderConfig(),
			DisableStacktrace: true,
			OutputPaths: []string{
				"stdout",
			},
		}
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		config.EncoderConfig.EncodeTime = devTimeEncoder
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

	return zap.Must(config.Build())
}

func devTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05"))
}
