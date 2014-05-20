package logo

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

//NewStdoutAppender creates an Appender that logs to standard output. Messages
//are formatted as follows: "[LEVEL]: [TIMESTAMP] [MESSAGE]"
func NewStdoutAppender() Appender {
	return &stdoutAppender{
		log.New(os.Stdout, "", 0), //use std logger for its buffering
	}
}

//NewLogstashAppender creates an Appender that emits logstash formatted messages
//to the specified writer.
func NewLogstashAppender(w io.Writer, v LogstashVersion, pretty bool) Appender {
	return &logstashAppender{
		log.New(w, "", 0),
		v,
		pretty,
	}
}

type stdoutAppender struct {
	logger *log.Logger
}

func (s *stdoutAppender) Write(message *LogMessage) {
	const layout = "2006-01-02T15:04:05.000Z"
	s.logger.Printf("%s: [%v] %s", message.Level, message.Timestamp.Format(layout), message.Message)

}

func (s *stdoutAppender) Close() { /* noop */
}

type LogstashVersion uint8

const (
	//logstash version 0
	LSV0 LogstashVersion = iota

	//logstash version 1
	LSV1
)

type logstashAppender struct {
	logger  *log.Logger
	version LogstashVersion
	pretty  bool
}

const logstashLayout = "2006-01-02T15:04:05.000000Z"

func (l *logstashAppender) Write(message *LogMessage) {
	var out interface{}
	switch l.version {
	case LSV0:
		msg := make(map[string]interface{})
		message.MDC["level"] = message.Level.String()
		msg["@fields"] = message.MDC
		msg["@message"] = message.Message
		msg["@timestamp"] = message.Timestamp.Format(logstashLayout)
		out = msg
	default: //use version 1
		msg := message.MDC
		msg["level"] = message.Level.String()
		msg["@message"] = message.Message
		msg["@timestamp"] = message.Timestamp.Format(logstashLayout)
		out = msg
	}

	var logData []byte
	if l.pretty {
		logData, _ = json.MarshalIndent(out, "", "  ")
	} else {
		logData, _ = json.Marshal(out)
	}

	l.logger.Print(string(logData))
}

func (l *logstashAppender) Close() { /* noop */
}
