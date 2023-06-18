package logger

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"Cryptogo/config"
)

type Logger struct {
	*zerolog.Logger
}

var (
	logger *Logger
	once   sync.Once
)

func GetLogger(cfg *config.Config) *Logger {
	once.Do(func() {
		loggingLevel := cfg.Server.LoggingLevel

		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		output.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		}
		output.FormatMessage = func(i interface{}) string {
			return fmt.Sprintf("%s", i)
		}
		output.FormatFieldName = func(i interface{}) string {
			return fmt.Sprintf("%s:", i)
		}
		output.FormatFieldValue = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("%s", i))
		}

		zeroLogger := zerolog.New(output).With().Timestamp().Logger()

		switch loggingLevel {
		case "debug":
			zerolog.SetGlobalLevel(zerolog.DebugLevel)
		case "error":
			zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		case "fatal":
			zerolog.SetGlobalLevel(zerolog.FatalLevel)
		case "panic":
			zerolog.SetGlobalLevel(zerolog.PanicLevel)
		default:
			zerolog.SetGlobalLevel(zerolog.NoLevel)
		}
		logger = &Logger{&zeroLogger}
	})

	return logger
}
