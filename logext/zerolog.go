package logext

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func ConsoleWriter() *zerolog.Logger {
	logger := log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
	return &logger
}
