package logger

import (
	"os"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/vinicius73/thecollector/pkg/support"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.DurationFieldInteger = true
}

func SetupLogger(cgf Config, level string, tags map[string]interface{}) {
	if level == "" {
		level = cgf.Level
	}

	zerolog.SetGlobalLevel(getLogLevel(level))

	log.Logger = buildBaseLogger(log.Logger, cgf.Format).
		With().
		Fields(tags).
		Logger()
}

func Logger(process string, tags map[string]interface{}) zerolog.Logger {
	builder := log.Logger.With()

	if process != "" {
		builder = builder.Str("process", process)
	}

	return builder.Fields(tags).Logger()
}

func getLogLevel(val string) zerolog.Level {
	if val == "" {
		val = support.GetEnv("LOG_LEVEL", "info")
	}

	level := strings.ToLower(val)

	switch level {
	case "debug":
		return zerolog.DebugLevel
	case "trace":
		return zerolog.TraceLevel
	default:
		return zerolog.InfoLevel
	}
}

func buildBaseLogger(logger zerolog.Logger, format string) zerolog.Logger {
	logger = logger.With().Str("app", "thecollectior").Logger()

	if format == "" {
		format = strings.ToLower(support.GetEnv("LOG_FORMAT", "json"))
	}

	switch format {
	case "json":
		return logger
	default:
		return logger.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
}
