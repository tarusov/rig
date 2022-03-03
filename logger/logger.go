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
		log *zerolog.Logger

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
func New(opts ...LoggerOption) *Logger {

	var log = &Logger{
		outputs:         []io.Writer{defaultOutput},
		level:           defaultLevel,
		format:          defaultFormat,
		timestampName:   zerolog.TimestampFieldName,
		timestampFormat: defaultTimeFormat,
	}

	for _, opt := range opts {
		opt(log)
	}

	if log.format == FormatText || log.format == FormatConsole {
		for i, o := range log.outputs {
			switch o.(type) {
			case *os.File:
				log.outputs[i] = &zerolog.ConsoleWriter{
					Out:        o,
					NoColor:    (log.format != FormatConsole),
					TimeFormat: log.timestampFormat,
				}
			default:
				continue
			}
		}
	}

	zerolog.TimestampFieldName = log.timestampName

	var lw = zerolog.MultiLevelWriter(log.outputs...)
	var l = zerolog.New(lw).With().Timestamp().Logger().Level(log.level)

	log.log = &l

	return log
}

// Debug implements Debug method for logger.
func (l *Logger) Debug(msg string) {
	l.log.Debug().Msg(msg)
}

// Info implements Info method for logger.
func (l *Logger) Info(msg string) {
	l.log.Info().Msg(msg)
}

// Warn implements Warn method for logger.
func (l *Logger) Warn(msg string) {
	l.log.Warn().Msg(msg)
}

// Error implements Error method for logger.
func (l *Logger) Error(msg string) {
	l.log.Error().Msg(msg)
}

// Fatal implements Fatal method for logger.
func (l *Logger) Fatal(msg string) {
	l.log.Fatal().Msg(msg)
}

// Debugf implements Debugf method for logger.
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.log.Debug().Msgf(format, args...)
}

// Infof implements Infof method for logger.
func (l *Logger) Infof(format string, args ...interface{}) {
	l.log.Info().Msgf(format, args...)
}

// Warnf implements Warnf method for logger.
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.log.Warn().Msgf(format, args...)
}

// Errorf implements Errorf method for logger.
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.log.Error().Msgf(format, args...)
}

// Fatalf implements Fatalf method for logger.
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.log.Fatal().Msgf(format, args...)
}

// WithField implements WithField method for logger.
func (l *Logger) WithField(name string, value interface{}) *Logger {
	var outzl = l.log.With().Interface(name, value).Logger()
	return &Logger{log: &outzl}
}

// Auxilary type for method WithFields (map string-interface).
type Fields map[string]interface{}

// WithFields implements WithFields method for logger.
func (l *Logger) WithFields(fields Fields) *Logger {

	var outzl = *l.log
	for k, v := range fields {
		outzl = outzl.With().Interface(k, v).Logger()
	}

	return &Logger{log: &outzl}
}

// WithErr implements WithErr method for logger.
func (l *Logger) WithErr(err error) *Logger {
	var outzl = l.log.With().Err(err).Logger()
	return &Logger{log: &outzl}
}

// globalLogger default logger.
var globalLogger = New()

// Global return registered global logger (or default).
func Global() *Logger {
	return globalLogger
}

// SetGlobal setup global logger.
func SetGlobal(logger *Logger) {
	globalLogger = logger
}
