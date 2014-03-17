package logo

import (
	"log"
	"os"
)

func NewStdoutAppender() Appender {
	return &stdoutAppender{
		log.New(os.Stdout, "", 0), //use logger for its buffering
	}
}

type stdoutAppender struct {
	logger *log.Logger
}

func (s *stdoutAppender) Write(message *LogMessage) {
	const layout = "2006-01-02T15:04:05.000Z"
	s.logger.Printf("%s: [%v] %s", message.Level, message.Timestamp.Format(layout), message.Message)
	
}

func (s *stdoutAppender) Close() {/* noop */}

















