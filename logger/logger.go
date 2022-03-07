// Package logger contains Logger implementation based on zerolog.
package logger

import (
	"io"
	"os"
	"time"

	"github.com/rs/zerolog"
)

// Format is output formatting option.
type Format string

// Format type enum.
const (
	FormatConsole Format = "console"
	FormatJSON    Format = "json"
	FormatText    Format = "text"
)

// Level is event's severnity level.
type Level string

// Level type enum.
const (
	LevelDebug   Level = "debug"
	LevelInfo    Level = "info"
	LevelWarning Level = "warn"
	LevelError   Level = "error"
	LevelFatal   Level = "fatal"
)

type (
	// Logger struct.
	Logger struct {
		*zerolog.Logger
	}

	// loggerOptions is auxilary constructor struct.
	loggerOptions struct {
		outputs         []io.Writer
		level           zerolog.Level
		format          Format
		timestampName   string
		timestampFormat string
	}
)

// Defaults.
const (
	defaultFormat     = FormatJSON
	defaultLevel      = zerolog.DebugLevel
	defaultTimeFormat = time.RFC3339
)

// defaultOutput is standart log output target.
var defaultOutput = os.Stderr

// New func create new logger instance.
func New(opts ...loggerOption) *Logger {

	var lo = &loggerOptions{
		outputs:         []io.Writer{defaultOutput},
		level:           defaultLevel,
		format:          defaultFormat,
		timestampName:   zerolog.TimestampFieldName,
		timestampFormat: defaultTimeFormat,
	}

	for _, opt := range opts {
		opt(lo)
	}

	if lo.format == FormatText || lo.format == FormatConsole {
		for i, o := range lo.outputs {
			switch o.(type) {
			case *os.File:
				lo.outputs[i] = &zerolog.ConsoleWriter{
					Out:        o,
					NoColor:    (lo.format != FormatConsole),
					TimeFormat: lo.timestampFormat,
				}
			default:
				continue
			}
		}
	}

	zerolog.TimestampFieldName = lo.timestampName

	var lw = zerolog.MultiLevelWriter(lo.outputs...)
	var zl = zerolog.New(lw).With().Timestamp().Logger().Level(lo.level)

	return &Logger{Logger: &zl}
}

// Debug implements Debug method for logger.
func (l *Logger) Debug(msg string) {
	l.Logger.Debug().Msg(msg)
}

// Info implements Info method for logger.
func (l *Logger) Info(msg string) {
	l.Logger.Info().Msg(msg)
}

// Warn implements Warn method for logger.
func (l *Logger) Warn(msg string) {
	l.Logger.Warn().Msg(msg)
}

// Error implements Error method for logger.
func (l *Logger) Error(msg string) {
	l.Logger.Error().Msg(msg)
}

// Fatal implements Fatal method for logger.
func (l *Logger) Fatal(msg string) {
	l.Logger.Fatal().Msg(msg)
}

// Debugf implements Debugf method for logger.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.Logger.Debug().Msgf(format, args...)
}

// Infof implements Infof method for logger.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.Logger.Info().Msgf(format, args...)
}

// Warnf implements Warnf method for logger.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.Logger.Warn().Msgf(format, args...)
}

// Errorf implements Errorf method for logger.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.Logger.Error().Msgf(format, args...)
}

// Fatalf implements Fatalf method for logger.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.Logger.Fatal().Msgf(format, args...)
}

// WithField implements WithField method for logger.
func (l *Logger) WithField(key string, value interface{}) *Logger {
	nk, nv := normalize(key, value)
	var outzl = l.Logger.With().Str(nk, nv).Logger()
	return &Logger{Logger: &outzl}
}

// Auxilary type for method WithFields (map string-interface).
type Fields map[string]interface{}

// WithFields implements WithFields method for logger.
func (l *Logger) WithFields(fields Fields) *Logger {

	var outzl = *l.Logger
	for k, v := range fields {
		nk, nv := normalize(k, v)
		outzl = outzl.With().Str(nk, nv).Logger()
	}

	return &Logger{Logger: &outzl}
}

// WithErr implements WithErr method for logger. Do not add error if no error pushed.
func (l *Logger) WithErr(err error) *Logger {
	if err != nil {
		var outzl = l.Logger.With().Err(err).Logger()
		return &Logger{Logger: &outzl}
	}
	return l
}
