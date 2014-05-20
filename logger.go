package logo

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
	"time"
)

type Level uint8

const (
	TRACE Level = iota
	INFO
	WARN
	ERROR
	PANIC
	NONE
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
	case PANIC:
		return "PANIC"
	default:
		return "INFO"
	}
}

//LogMessage carries the formatted string that was logged along with contextual
//meta-data
type LogMessage struct {
	Level     Level
	Timestamp time.Time
	Message   string
	MDC       map[string]string
}

//Context is the base context that is inherited by all loggers created. Any
//attribute added to this context will be included in the MDC of every log
//emitted
var Context *MDC

//LogLevel is the current log level. The logger will ignore any message sent to
//it that are lower than the specified level. The default level is 0 (TRACE)
var LogLevel Level

func init() {
	Context = NewMDC()
}

//NewLogger creates a new Logger that overlays the provided context onto the
//base context specified by Context.
func NewLogger(ctx map[string]string) *Logger {
	mdc := Context.snapshot()
	for key, value := range ctx {
		mdc[key] = value
	}
	return &Logger{MDCFromMap(mdc), 0}
}

//Logger is the receiver of log messages. All log messages sent to a logger will
//contain the MDC of the logger.
type Logger struct {
	MDC *MDC
	//StackDepth is used to adjust the file name lookups for trace logs.
	StackDepth uint8
}

//NewLogger returns a new logger building on the context of this logger
func (l *Logger) NewLogger(ctx map[string]string) *Logger {
	mdc := l.MDC.snapshot()
	for key, value := range ctx {
		mdc[key] = value
	}
	return &Logger{MDCFromMap(mdc), 0}
}

func (l *Logger) Tracef(arg0 string, args ...interface{}) {
	if LogLevel > TRACE {
		return
	}
	mdc := l.MDC.snapshot()
	mdc["file"] = getFileStr(2 + int(l.StackDepth))
	logMessage(TRACE, mdc, fmt.Sprintf(arg0, args...))
}

func (l *Logger) Infof(arg0 string, args ...interface{}) {
	if LogLevel > INFO {
		return
	}
	logMessage(INFO, l.MDC.snapshot(), fmt.Sprintf(arg0, args...))
}

func (l *Logger) Warnf(arg0 string, args ...interface{}) {
	if LogLevel > WARN {
		return
	}
	logMessage(WARN, l.MDC.snapshot(), fmt.Sprintf(arg0, args...))
}

func (l *Logger) Errorf(arg0 string, args ...interface{}) {
	if LogLevel > ERROR {
		return
	}
	logMessage(ERROR, l.MDC.snapshot(), fmt.Sprintf(arg0, args...))
}

func (l *Logger) Panicf(arg0 string, args ...interface{}) {
	if LogLevel > PANIC {
		return
	}

	//get a stack trace
	stack := make([]byte, 1024)
	size := runtime.Stack(stack, true)
	stackStr := string(stack[:size])

	//add diagnostic info to context
	mdc := l.MDC.snapshot()
	mdc["file"] = getFileStr(2 + int(l.StackDepth))
	mdc["stack_trace"] = stackStr

	message := fmt.Sprintf(arg0, args...)
	logMessage(PANIC, mdc, message)
	panic(errors.New(message))
}

func getFileStr(skip int) string {
	if _, file, line, ok := runtime.Caller(skip); ok {
		idx := strings.LastIndex(file, "/")
		return fmt.Sprintf("%s:%d", file[idx+1:], line)
	} else {
		return "???:0"
	}
}
