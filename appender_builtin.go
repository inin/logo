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
	s.logger.Printf("%s: [%v] %s", message.Level, message.Timestamp, message.Message)
	
}

func (s *stdoutAppender) Close() {/* noop */}

















