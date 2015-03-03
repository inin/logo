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
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
	PANIC
	NONE
)

func (l Level) String() string {
	switch l {
	case TRACE:
		return "TRACE"
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
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
func (this *Logger) NewLogger(ctx map[string]string) *Logger {
	mdc := this.MDC.snapshot()
	for key, value := range ctx {
		mdc[key] = value
	}
	return &Logger{MDCFromMap(mdc), 0}
}

func (this* Logger) isLevelEnabled(level Level) bool {
	return LogLevel <= level
}

func (this* Logger) IsTraceEnabled() bool {
	return this.isLevelEnabled(TRACE)
}

func (this* Logger) IsDebugEnabled() bool {
	return this.isLevelEnabled(DEBUG)
}

func (this* Logger) IsInfoEnabled() bool {
	return this.isLevelEnabled(INFO)
}

func (this* Logger) IsWarnEnabled() bool {
	return this.isLevelEnabled(WARN)
}

func (this* Logger) IsErrorEnabled() bool {
	return this.isLevelEnabled(ERROR)
}

func (this* Logger) IsFatalEnabled() bool {
	return this.isLevelEnabled(FATAL)
}

func (this* Logger) IsPanicEnabled() bool {
	return this.isLevelEnabled(DEBUG)
}

func (this* Logger) logAt(level Level, attributes map[string]string, arg0 string, args ...interface{}) {
	// check the level first
	if !this.isLevelEnabled(level) {
		return
	}

	// clone the logger context
	mdc := this.MDC.snapshot()

	// attach the additional attributes, if any
	if attributes != nil {
		for key, value := range attributes {
			mdc[key] = value
		}
	}

	// FIXME: getting the file locatio is tricky since the call stack depths are different
	// depending on how the logger has been invoked. For example, any higher-level wrapper functions
	// will change the stack depth at which the interesting call frame is located. We can probably pass the expected
	// depth, but this is adding unnecessary complexity
	//mdc["file"] = getFileStr(4 + int(this.StackDepth))
	
	//add diagnostic info to context
	if level == FATAL || level == PANIC {
		//get a stack trace
		stack := make([]byte, 1024)
		size := runtime.Stack(stack, true)
		stackStr := string(stack[:size])
		mdc["stack_trace"] = stackStr
	}

	message := fmt.Sprintf(arg0, args...)

	logMessage(level, mdc, message)
	
	// FIXME: while this is convenient, it combines tracing with control flow and has side effects, which should probably
	// be avoided. For example, is the log level is set to NONE, then panic will not be invoked and the control flow
	// will not be affected. It may be better to move the panic call out of here, but there is probably code out there
	// that relies on this behavior.
	if level == PANIC {
		panic(errors.New(message))
	}
}

func (this *Logger) Tracef(arg0 string, args ...interface{}) {
	this.logAt(TRACE, nil, arg0, args...)
}

func (this *Logger) Debugf(arg0 string, args ...interface{}) {
	this.logAt(DEBUG, nil, arg0, args...)
}

func (this *Logger) Infof(arg0 string, args ...interface{}) {
	this.logAt(INFO, nil, arg0, args...)
}

func (this *Logger) Warnf(arg0 string, args ...interface{}) {
	this.logAt(WARN, nil, arg0, args...)
}

func (this *Logger) Errorf(arg0 string, args ...interface{}) {
	this.logAt(ERROR, nil, arg0, args...)
}

func (this *Logger) Fatalf(arg0 string, args ...interface{}) {
	this.logAt(FATAL, nil, arg0, args...)
}

func (this *Logger) Panicf(arg0 string, args ...interface{}) {
	this.logAt(PANIC, nil, arg0, args...)
}

func (this *Logger) ContextTracef(ctx map[string]string, arg0 string, args ...interface{}) {
	this.logAt(TRACE, ctx, arg0, args...)
}

func (this *Logger) ContextDebugf(ctx map[string]string, arg0 string, args ...interface{}) {
	this.logAt(DEBUG, ctx, arg0, args...)
}

func (this *Logger) ContextInfof(ctx map[string]string, arg0 string, args ...interface{}) {
	this.logAt(INFO, ctx, arg0, args...)
}

func (this *Logger) ContextWarnf(ctx map[string]string, arg0 string, args ...interface{}) {
	this.logAt(WARN, ctx, arg0, args...)
}

func (this *Logger) ContextErrorf(ctx map[string]string, arg0 string, args ...interface{}) {
	this.logAt(ERROR, ctx, arg0, args...)
}

func (this *Logger) ContextFatalf(ctx map[string]string, arg0 string, args ...interface{}) {
	this.logAt(FATAL, ctx, arg0, args...)
}

func (this *Logger) ContextPanicf(ctx map[string]string, arg0 string, args ...interface{}) {
	this.logAt(PANIC, ctx, arg0, args...)
}

func getFileStr(skip int) string {
	if _, file, line, ok := runtime.Caller(skip); ok {
		idx := strings.LastIndex(file, "/")
		return fmt.Sprintf("%s:%d", file[idx+1:], line)
	} else {
		return "???:0"
	}
}
