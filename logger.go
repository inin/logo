package logo

import (
	"time"
	"fmt"
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

type Logger struct {
}

func (l *Logger) Tracef(arg0 string, args ...interface{}) {
	logMessage(TRACE, nil, fmt.Sprintf(arg0, args...))
}

func (l *Logger) Infof(arg0 string, args ...interface{}) {
	logMessage(INFO, nil, fmt.Sprintf(arg0, args...))
}

func (l *Logger) Warnf(arg0 string, args ...interface{}) {
	logMessage(WARN, nil, fmt.Sprintf(arg0, args...))
}

func (l *Logger) Errorf(arg0 string, args ...interface{}) {
	logMessage(ERROR, nil, fmt.Sprintf(arg0, args...))
}



