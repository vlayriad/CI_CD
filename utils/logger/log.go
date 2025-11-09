package logger

import (
	"cmp"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func baseConsoleWriter(loc *time.Location) zerolog.ConsoleWriter {
	return zerolog.ConsoleWriter{
		Out:          os.Stdout,
		TimeFormat:   zerolog.TimeFieldFormat,
		TimeLocation: loc,
		FormatTimestamp: func(i interface{}) string {
			str := fmt.Sprintf("%v", i)
			if len(str) >= 19 {
				str = str[:10] + " " + str[11:]
			}
			return fmt.Sprintf("[ %v ]", str)
		},
	}
}

func levelMode(lvlMode string) zerolog.Level {
	switch lvlMode {
	case "debug": // 0
		return zerolog.DebugLevel
	case "info": // 1
		return zerolog.InfoLevel
	case "warn": // 2
		return zerolog.WarnLevel
	case "error": // 3
		return zerolog.ErrorLevel
	case "fatal": // 4
		return zerolog.FatalLevel
	case "panic": // 5
		return zerolog.PanicLevel
	default:
		return zerolog.InfoLevel // 1
	}
}

func InitDefault() {
	loc, _ := time.LoadLocation("Asia/Jakarta")

	zerolog.TimeFieldFormat = "2006-01-02 15:04:05Z07:00"
	zerolog.TimestampFieldName = "timestamp"

	consoleWriter := baseConsoleWriter(loc)

	log.Logger = zerolog.New(consoleWriter).
		With().
		Timestamp().
		Caller().
		Logger()
}

func Init(timezone, appEnvironment, appDebug string) {
	loc, err := time.LoadLocation(cmp.Or(timezone, "Asia/Jakarta"))
	if err != nil {
		log.Error().Err(err).Msg("failed to load timezone, fallback to Asia/Jakarta")
		loc, _ = time.LoadLocation("Asia/Jakarta")
	}

	consoleWriter := baseConsoleWriter(loc)

	var multiWriter io.Writer = consoleWriter
	if appEnvironment == "production" {
		logFile, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Error().Err(err).Msg("failed to open log file, fallback to console only")
		} else {
			multiWriter = io.MultiWriter(consoleWriter, logFile)
		}

		zerolog.SetGlobalLevel(levelMode(appDebug))
	}

	logger := zerolog.New(multiWriter).
		With().
		Timestamp().
		Caller().
		Logger()

	log.Logger = logger
}
