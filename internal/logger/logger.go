package logger

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	"os"
	"time"
)

func New(logLevel string, logPath string) (*zerolog.Logger, *os.File) {
	var logger zerolog.Logger

	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zerolog.TimeFieldFormat = time.RFC3339Nano

	parsedLevel, err := zerolog.ParseLevel(logLevel)
	if err != nil {
		parsedLevel = zerolog.InfoLevel
	}

	if logPath != "" {
		file, err := os.OpenFile(
			logPath,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY,
			0664,
		)

		if err != nil {
			panic(err)
		}

		logger = zerolog.New(file).Level(parsedLevel).With().Timestamp().Logger()

		return &logger, file

	} else {
		logger = zerolog.New(os.Stdout).Level(parsedLevel).With().Timestamp().Logger()

		return &logger, nil
	}
}
