package shell

import "time"

type Logger interface {
	LogDuration() (time.Time, time.Time)

	doLog(level int, args ...any)
	doLogln(level int, args ...any)
	doLogf(level int, format string, args ...any)
	doLogfln(level int, format string, args ...any)

	Log(args ...any)
	Logln(args ...any)
	Logf(format string, args ...any)
	Logfln(format string, args ...any)

	Error(args ...any)
	Errorln(args ...any)
	Errorf(format string, args ...any)
	Errorfln(format string, args ...any)

	Warn(args ...any)
	Warnln(args ...any)
	Warnf(format string, args ...any)
	Warnfln(format string, args ...any)

	Info(args ...any)
	Infoln(args ...any)
	Infof(format string, args ...any)
	Infodln(format string, args ...any)

	Debug(args ...any)
	Debugln(args ...any)
	Debugf(format string, args ...any)
	Debugfln(format string, args ...any)

	Trace(args ...any)
	Traceln(args ...any)
	Tracef(format string, args ...any)
	Tracefln(format string, args ...any)
}
