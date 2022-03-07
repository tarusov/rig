package logger_test

import (
	"bytes"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/tarusov/rig/logger"
)

func TestLoggerLevelsDebug(t *testing.T) {

	var (
		buf = bytes.NewBuffer(make([]byte, 0))
		log = logger.New(
			logger.WithLoggingTimestampField("ts"),
			logger.WithLoggingTimestampFormat(time.Kitchen),
			logger.WithLoggingFormat(logger.FormatText),
			logger.WithLoggingLevel(logger.LevelDebug), // Must log all
			logger.WithLoggingOutput(buf),
		)
		conds = []string{"dbg", "dbgf", "inf", "inff", "wrn", "wrnf", "err", "errf"}
	)

	for _, cond := range conds {
		t.Run(fmt.Sprintf("log %s", cond), func(t *testing.T) {
			switch cond {
			case "dbg":
				log.Debug(cond)
			case "dbgf":
				log.Debugf(cond)
			case "inf":
				log.Info(cond)
			case "inff":
				log.Infof(cond)
			case "wrn":
				log.Warn(cond)
			case "wrnf":
				log.Warnf(cond)
			case "err":
				log.Error(cond)
			case "errf":
				log.Errorf(cond)
			}

			require.Containsf(t, buf.String(), cond, "TestLoggerLevelsDebug: message with level %q not logged", cond)
			buf.Reset()
		})
	} // For
}

func TestLoggerLevelsFatal(t *testing.T) {

	var (
		buf = bytes.NewBuffer(make([]byte, 0))
		log = logger.New(
			logger.WithLoggingTimestampField("ts"),
			logger.WithLoggingTimestampFormat(time.Kitchen),
			logger.WithLoggingFormat(logger.FormatText),
			logger.WithLoggingLevel(logger.LevelFatal), // Must skip all
			logger.WithLoggingOutput(buf),
		)
		conds = []string{"dbg", "dbgf", "inf", "inff", "wrn", "wrnf", "err", "errf"}
	)

	for _, cond := range conds {
		t.Run(fmt.Sprintf("log %s", cond), func(t *testing.T) {
			switch cond {
			case "dbg":
				log.Debug(cond)
			case "dbgf":
				log.Debugf(cond)
			case "inf":
				log.Info(cond)
			case "inff":
				log.Infof(cond)
			case "wrn":
				log.Warn(cond)
			case "wrnf":
				log.Warnf(cond)
			case "err":
				log.Error(cond)
			case "errf":
				log.Errorf(cond)
			}

			require.Emptyf(t, buf.String(), "TestLoggerLevelsFatal: message with level %q logged", cond)
			buf.Reset()
		})
	} // For
}

func TestLoggerFields(t *testing.T) {

	var (
		buf = bytes.NewBuffer(make([]byte, 0))
		log = logger.New(
			logger.WithLoggingTimestampField("ts"),
			logger.WithLoggingTimestampFormat(time.Kitchen),
			logger.WithLoggingFormat(logger.FormatText),
			logger.WithLoggingLevel(logger.LevelInfo),
			logger.WithLoggingOutput(buf),
		)
	)

	log.WithField("field_one", "1").Info("msg")
	require.Containsf(t, buf.String(), "field_one", "TestLoggerFields: message field not logged")
	buf.Reset()

	log.WithFields(logger.Fields{
		"field_one": "1",
		"field_two": "2",
	}).Info("msg")
	require.Containsf(t, buf.String(), "field_one", "TestLoggerFields: message fields not logged")
	require.Containsf(t, buf.String(), "field_two", "TestLoggerFields: message fields not logged")
	buf.Reset()

	log.WithErr(errors.New("some_err")).Info("msg")
	require.Containsf(t, buf.String(), "some_err", "TestLoggerFields: message error field not logged")
}
