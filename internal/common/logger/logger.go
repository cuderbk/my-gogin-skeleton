package logger

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"my-gogin-skeleton/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func New(cfg config.LogConfig) *zap.SugaredLogger {
	if Logger != nil {
		return Logger
	}

	var level zapcore.Level
	if err := level.UnmarshalText([]byte(strings.ToLower(cfg.Level))); err != nil {
		level = zapcore.InfoLevel
	}

	var encoderCfg zapcore.EncoderConfig
	if cfg.Pretty {
		encoderCfg = zap.NewDevelopmentEncoderConfig()
	} else {
		encoderCfg = zap.NewProductionEncoderConfig()
	}
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	var encoder zapcore.Encoder
	if cfg.Pretty {
		encoder = zapcore.NewConsoleEncoder(encoderCfg)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderCfg)
	}

	var output zapcore.WriteSyncer
	if strings.ToLower(cfg.Output) == "file" {
		logFile := cfg.Filename
		if logFile == "" {
			log.Fatalf("log output is 'file' but no 'filename' is provided")
		}

		logDir := filepath.Dir(logFile)

		if err := os.MkdirAll(logDir, 0755); err != nil {
			log.Fatalf("failed to create log dir: %v", err)
		}
		f, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatalf("failed to open log file: %v", err)
		}
		output = zapcore.AddSync(f)
	} else {
		output = zapcore.AddSync(os.Stdout)
	}

	core := zapcore.NewCore(encoder, output, level)
	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1)) // skip wrapper
	Logger = logger.Sugar()
	Logger.Infof("logger initialized [level=%s] [pretty=%v] [output=%s]", cfg.Level, cfg.Pretty, cfg.Output)

	return Logger
}
