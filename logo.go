package logo

import (
	"time"
)

var appenders = newAppenderList()
var msgCh = make(chan *LogMessage, 100)
var killCh = make(chan bool)
var doneCh = make(chan bool)

func init() {
	go listen()
}

func Close() {
	killCh <- true       //stop the listener
	<-doneCh             //wait for it to flush and complete
	appenders.closeAll() //close the appenders
}

func listen() {
	//keep processing messages until a kill is signaled
	for run := true; run; {
		select {
		case msg := <-msgCh:
			appenders.writeAll(msg)
		case <-killCh:
			run = false
		}
	}

	//try to flush the contents of the log channel
	for flush := true; flush; {
		select {
		case msg := <-msgCh:
			appenders.writeAll(msg)
		case <-time.Tick(3 * time.Second): //too many messages, abort
			flush = false
		default: //no more messages, abort
			flush = false
		}
	}

	doneCh <- true
}

func AddAppender(appender Appender) {
	appenders.add(appender)
}

func logMessage(level Level, mdc map[string]string, msg string) {
	message := &LogMessage{
		level,
		time.Now(),
		msg,
		mdc,
	}

	msgCh <- message
}
