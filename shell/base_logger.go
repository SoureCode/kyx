package shell

import (
	"fmt"
	"sync/atomic"
	"time"
)

var lastTimestamp time.Time
var firstTimestamp time.Time
var lastMessageHadNewline atomic.Bool
var loggerInstance *BaseLogger

var consoleLogHandlerInstance *ConsoleLogHandler

func GetConsoleLogHandler() *ConsoleLogHandler {
	if consoleLogHandlerInstance == nil {
		consoleLogHandlerInstance = NewConsoleLogHandler()
	}

	return consoleLogHandlerInstance
}

func GetLogger() *BaseLogger {
	if loggerInstance == nil {
		loggerInstance = newLogger()
	}

	return loggerInstance
}

type BaseLogger struct {
	handlers []LogHandler
}

func newLogger() *BaseLogger {
	logger := &BaseLogger{
		handlers: make([]LogHandler, 0),
	}

	lastMessageHadNewline.Store(true)

	logger.AddHandler(GetConsoleLogHandler())

	return logger
}

func (l *BaseLogger) AddHandler(handler LogHandler) {
	l.handlers = append(l.handlers, handler)
}

func (l *BaseLogger) ClearHandlers() {
	l.handlers = make([]LogHandler, 0)
}

func (l *BaseLogger) now() time.Time {
	now := time.Now()

	if firstTimestamp.IsZero() {
		firstTimestamp = now
	}

	return now
}
func (l *BaseLogger) getTimestamp(t time.Time) string {
	return t.Format(time.RFC3339)
}

func (l *BaseLogger) getDelta(t time.Time) string {
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

func (l *BaseLogger) LogDuration() {
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

func (l *BaseLogger) doLog(level int, args ...any) {
	if lastMessageHadNewline.Load() {
		now := l.now()
		timestamp := l.getTimestamp(now)
		delta := l.getDelta(now)

		prefix := []any{wrapInBrackets(timestamp) + wrapInBrackets(delta) + " "}
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

	for _, handler := range l.handlers {
		if handler != nil {
			handler.Log(level, args...)
		}
	}
}

func (l *BaseLogger) doLogln(level int, args ...any) {
	args = append(args, "\n")

	l.doLog(level, args...)
}

func (l *BaseLogger) doLogf(level int, format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	l.doLog(level, message)
}

func (l *BaseLogger) doLogfln(level int, format string, args ...any) {
	message := fmt.Sprintf(format, args...)
	l.doLogln(level, message)
}

func (l *BaseLogger) Log(args ...any) {
	l.doLog(0, args...)
}

func (l *BaseLogger) Logln(args ...any) {
	l.doLogln(0, args...)
}

func (l *BaseLogger) Logf(format string, args ...any) {
	l.doLogf(0, format, args...)
}

func (l *BaseLogger) Logfln(format string, args ...any) {
	l.doLogfln(0, format, args...)
}

func (l *BaseLogger) Error(args ...any) {
	l.doLog(1, args...)
}
func (l *BaseLogger) Errorln(args ...any) {
	l.doLogln(1, args...)
}
func (l *BaseLogger) Errorf(format string, args ...any) {
	l.doLogf(1, format, args...)
}
func (l *BaseLogger) Errorfln(format string, args ...any) {
	l.doLogfln(1, format, args...)
}

func (l *BaseLogger) Warn(args ...any) {
	l.doLog(2, args...)
}
func (l *BaseLogger) Warnln(args ...any) {
	l.doLogln(2, args...)
}
func (l *BaseLogger) Warnf(format string, args ...any) {
	l.doLogf(2, format, args...)
}
func (l *BaseLogger) Warnfln(format string, args ...any) {
	l.doLogfln(2, format, args...)
}

func (l *BaseLogger) Info(args ...any) {
	l.doLog(3, args...)
}
func (l *BaseLogger) Infoln(args ...any) {
	l.doLogln(3, args...)
}
func (l *BaseLogger) Infof(format string, args ...any) {
	l.doLogf(3, format, args...)
}
func (l *BaseLogger) Infodln(format string, args ...any) {
	l.doLogfln(3, format, args...)
}
func (l *BaseLogger) Debug(args ...any) {
	l.doLog(4, args...)
}
func (l *BaseLogger) Debugln(args ...any) {
	l.doLogln(4, args...)
}
func (l *BaseLogger) Debugf(format string, args ...any) {
	l.doLogf(4, format, args...)
}
func (l *BaseLogger) Debugfln(format string, args ...any) {
	l.doLogfln(4, format, args...)
}
func (l *BaseLogger) Trace(args ...any) {
	l.doLog(5, args...)
}
func (l *BaseLogger) Traceln(args ...any) {
	l.doLogln(5, args...)
}
func (l *BaseLogger) Tracef(format string, args ...any) {
	l.doLogf(5, format, args...)
}
func (l *BaseLogger) Tracefln(format string, args ...any) {
	l.doLogfln(5, format, args...)
}
