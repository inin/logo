package logo

import (
	"fmt"
	"time"
)

type Level uint8

const (
	TRACE Level = iota
	INFO
	WARN
	ERROR
)

func (l Level) String() string {
	switch l {
	case TRACE:
		return "TRACE"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "INFO"
	}
}

type LogMessage struct {
	Level     Level
	Timestamp time.Time
	Message   string
	MDC       map[string]string
}

var Context *MDC

func init() {
	Context = NewMDC()
}

func NewLogger(ctx map[string]string) *Logger {
	mdc := Context.snapshot()
	for key, value := range ctx {
		mdc[key] = value
	}
	return &Logger{MDCFromMap(mdc)}
}

type Logger struct {
	MDC *MDC
}

//NewLogger returns a new logger building on the context of this logger
func (l *Logger) NewLogger(ctx map[string]string) *Logger {
	mdc := l.MDC.snapshot()
	for key, value := range ctx {
		mdc[key] = value
	}
	return &Logger{MDCFromMap(mdc)}
}

func (l *Logger) Tracef(arg0 string, args ...interface{}) {
	logMessage(TRACE, l.MDC.snapshot(), fmt.Sprintf(arg0, args...))
}

func (l *Logger) Infof(arg0 string, args ...interface{}) {
	logMessage(INFO, l.MDC.snapshot(), fmt.Sprintf(arg0, args...))
}

func (l *Logger) Warnf(arg0 string, args ...interface{}) {
	logMessage(WARN, l.MDC.snapshot(), fmt.Sprintf(arg0, args...))
}

func (l *Logger) Errorf(arg0 string, args ...interface{}) {
	logMessage(ERROR, l.MDC.snapshot(), fmt.Sprintf(arg0, args...))
}
