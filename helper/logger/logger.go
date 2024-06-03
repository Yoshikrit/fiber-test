package logger

import (
	"github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

func init() {
    zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	// zerolog.SetGlobalLevel(zerolog.DebugLevel)
}

func Info(msg string, keysAndValues ...interface{}) {
	log.Info().Fields(keysAndValues).Msg(msg)
}

func Debug(msg string, keysAndValues ...interface{}) {
	log.Debug().Fields(keysAndValues).Msg(msg)
}

func Error(msg interface{}, keysAndValues ...interface{}) {
	switch v := msg.(type) {
	case error:
		log.Error().Fields(keysAndValues).Msg(v.Error())
	case string:
		log.Error().Fields(keysAndValues).Msg(v)
	}
}
