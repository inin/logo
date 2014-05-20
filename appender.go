package logo

import (
	"sync"
)

//Appender is an interface that allows direction of log messages to a variety of
//destinations. Appenders can be added via the AddAppender function.
type Appender interface {
	Write(message *LogMessage)
	Close()
}

type writer struct {
	appender Appender
	msgCh    chan *LogMessage
}

func (w *writer) listen(wg *sync.WaitGroup) {
	defer wg.Done()
	for msg := range w.msgCh {
		w.appender.Write(msg)
	}
	w.appender.Close()
}

func newAppenderList() *appenderList {
	return &appenderList{
		make([]writer, 0, 10),
		sync.RWMutex{},
		sync.WaitGroup{},
	}
}

type appenderList struct {
	writers []writer
	mu      sync.RWMutex
	wg      sync.WaitGroup
}

func (a *appenderList) add(appender Appender) {
	w := writer{
		appender,
		make(chan *LogMessage, 100),
	}

	//increment the wait group count and start the listener
	a.wg.Add(1)
	go w.listen(&a.wg)

	//aquire a writer lock on the writers
	a.mu.Lock()
	defer a.mu.Unlock()

	//only append if the list exists
	if a.writers != nil {
		a.writers = append(a.writers, w)
	}
}

func (a *appenderList) writeAll(msg *LogMessage) {
	//aquire a reader lock on the writers
	a.mu.RLock()
	defer a.mu.RUnlock()

	for _, writer := range a.writers {
		writer.msgCh <- msg
	}
}

func (a *appenderList) closeAll() {
	//aquire a reader lock on the writers
	a.mu.RLock()
	defer a.mu.RUnlock()

	//prevent any more log messages from coming to these writers
	writers := a.writers
	a.writers = nil

	//close the message channel on each writer
	for _, writer := range writers {
		close(writer.msgCh)
	}

	//wait for writers routines to stop
	a.wg.Wait()
}
