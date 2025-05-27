package shell

import (
	"fmt"
	"sync/atomic"
	"time"
)

// Log
// Logln
// Logf
// Logfln

var lastTimestamp time.Time
var firstTimestamp time.Time
var loggerInstance *Logger
var lastMessageHadNewline atomic.Bool

func GetLogger() *Logger {
	if loggerInstance == nil {
		loggerInstance = newLogger()
	}

	return loggerInstance
}

type Logger struct {
	scope string
}

func (l *Logger) WithScope(scope string) *Logger {
	newLogger := &Logger{
		scope: scope,
	}

	return newLogger
}

func newLogger() *Logger {
	logger := &Logger{
		scope: "",
	}

	lastMessageHadNewline.Store(true)

	return logger
}

func (l *Logger) now() time.Time {
	now := time.Now()

	if firstTimestamp.IsZero() {
		firstTimestamp = now
	}

	return now
}
func (l *Logger) getTimestamp(t time.Time) string {
	return t.Format(time.RFC3339)
}

func (l *Logger) getDelta(t time.Time) string {
	if lastTimestamp.IsZero() {
		lastTimestamp = t
		return "+00:00.000"
	}

	delta := t.Sub(lastTimestamp)
	lastTimestamp = t

	totalMs := delta.Milliseconds()
	minutes := totalMs / 60000
	seconds := (totalMs % 60000) / 1000
	milliseconds := totalMs % 1000

	return fmt.Sprintf("+%02d:%02d.%03d", minutes, seconds, milliseconds)
}

func (l *Logger) LogDuration() {
	now := l.now()
	d := now.Sub(firstTimestamp)

	l.Logln("Total duration: ", d.String())
}

func wrapInBrackets(s string) string {
	if s == "" {
		return ""
	}

	return "[" + s + "]"
}

func (l *Logger) Log(args ...any) {
	if lastMessageHadNewline.Load() {
		now := l.now()
		timestamp := l.getTimestamp(now)
		delta := l.getDelta(now)

		prefix := []any{wrapInBrackets(l.scope) + wrapInBrackets(timestamp) + wrapInBrackets(delta) + " "}
		args = append(prefix, args...)
	}

	if len(args) > 0 {
		if str, ok := args[len(args)-1].(string); ok && str[len(str)-1] == '\n' {
			lastMessageHadNewline.Store(true)
		} else {
			lastMessageHadNewline.Store(false)
		}
	} else {
		lastMessageHadNewline.Store(false)
	}

	fmt.Print(args...)
}

func (l *Logger) Logln(args ...any) {
	args = append(args, "\n")

	l.Log(args...)
}

func (l *Logger) Logf(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	l.Log(message)
}

func (l *Logger) Logfln(format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	l.Logln(message)
}
