package logger

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/rs/zerolog"
)

type (
	// SentryNotifier implements sentry notifier.
	SentryNotifier struct {
		hub *sentry.Hub

		level      zerolog.Level
		timeout    time.Duration
		env        string
		release    string
		stacktrace bool
	}
)

// NewSentryNotifier create new sentry notifier instance.
func NewSentryNotifier(dsn string, opts ...SentryOption) (*SentryNotifier, error) {

	if dsn == "" {
		return nil, errors.New("sentry dsn is empty")
	}

	var sn = &SentryNotifier{
		level:   defaultLevel,
		timeout: 5 * time.Second,
	}

	for _, opt := range opts {
		opt(sn)
	}

	var client, err = sentry.NewClient(sentry.ClientOptions{
		Dsn:              dsn,
		AttachStacktrace: sn.stacktrace,
		Debug:            sn.level == zerolog.DebugLevel,
		Environment:      sn.env,
		Release:          sn.release,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to init sentry client: %w", err)
	}

	var scope = sentry.NewScope()
	var hub = sentry.NewHub(client, scope)

	sn.hub = hub

	return sn, nil
}

// Write is implemens io.Writer Write method.
func (sn *SentryNotifier) Write(p []byte) (int, error) {
	sn.capture(string(p), sn.level, nil)
	return len(p), nil
}

// WriteLevel is implements zerolog.LevelWriter method.
func (sn *SentryNotifier) WriteLevel(zl zerolog.Level, p []byte) (int, error) {

	var count = len(p)

	if zl < sn.level {
		return count, nil
	}

	var extra Fields
	var err = json.Unmarshal(p, &extra)
	if err != nil {
		return 0, fmt.Errorf("unmarshal message: %w", err)
	}

	var msgField, msg = getErrMessage(extra)

	delete(extra, msgField)
	delete(extra, zerolog.MessageFieldName)
	delete(extra, zerolog.ErrorFieldName)
	delete(extra, zerolog.LevelFieldName)
	delete(extra, zerolog.TimestampFieldName)

	sn.capture(msg, zl, extra)

	return count, nil
}

// capture sends final message to server.
func (sn *SentryNotifier) capture(msg string, level zerolog.Level, extra Fields) {

	var e = sentry.NewEvent()
	var stacktrace *sentry.Stacktrace
	if sn.stacktrace {
		stacktrace = getStacktrace(msg)
	}

	e.Message = msg
	e.Level = sentry.Level(level.String())
	e.Timestamp = time.Now()
	e.Exception = []sentry.Exception{{
		Value:      msg,
		Stacktrace: stacktrace,
	}}

	_ = sn.hub.CaptureEvent(e)
}

// possible err message fields.
const (
	logFieldErr = "err"
	logFieldMsg = "message"
)

// getErrMessage from error extra data.
func getErrMessage(extra Fields) (field string, msg string) {
	var ok bool
	if msg, ok = extra[logFieldErr].(string); ok {
		return logFieldErr, msg
	}

	if msg, ok = extra[logFieldMsg].(string); ok {
		return logFieldMsg, msg
	}

	return logFieldMsg, "undefined error"
}

// getStacktrace extract stacktrace from error.
func getStacktrace(msg string) *sentry.Stacktrace {

	var stacktrace = sentry.ExtractStacktrace(errors.New(msg))
	if stacktrace == nil {
		stacktrace = sentry.NewStacktrace()
	}
	if stacktrace == nil {
		return nil
	}

	var frames = make([]sentry.Frame, 0, len(stacktrace.Frames))
	for _, frame := range stacktrace.Frames {
		// Skip tracing into logger.
		if strings.HasPrefix(frame.Module, "github.com/rs/zerolog") ||
			strings.HasSuffix(frame.Filename, "logger.go") ||
			strings.HasSuffix(frame.Filename, "sentry.go") {
			continue
		}

		frames = append(frames, frame)
	}
	stacktrace.Frames = frames

	return stacktrace
}

// Close implements io.Close method for sentry notifier.
func (sn *SentryNotifier) Close() error {
	if !sn.hub.Flush(sn.timeout) {
		return errors.New("sentry flush timeout")
	}
	return nil
}
