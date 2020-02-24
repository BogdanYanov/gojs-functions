package jstime

import (
	"sync"
	"time"
)

type jsFields struct {
	duration time.Duration
	stopCh   chan struct{}
	mu       *sync.Mutex
	isClosed bool
}

func newJSFields(millis uint32) *jsFields {
	return &jsFields{
		duration: time.Duration(millis) * time.Millisecond,
		stopCh:   make(chan struct{}),
		mu:       &sync.Mutex{},
		isClosed: false,
	}
}

func (jsF *jsFields) isChannelClosed() bool {
	jsF.mu.Lock()
	defer jsF.mu.Unlock()
	return jsF.isClosed
}

func (jsF *jsFields) stopExecute() {
	if !jsF.isChannelClosed() {
		jsF.stopCh <- struct{}{}
		jsF.closeChannel()
	}
}

func (jsF *jsFields) closeChannel() {
	jsF.mu.Lock()
	close(jsF.stopCh)
	jsF.isClosed = true
	jsF.mu.Unlock()
}

// SetInterval called function 'f' every 'millis' milliseconds. It also returns a function
// that stops the function call 'f'
func SetInterval(f func(), millis uint32, wg *sync.WaitGroup) func() {
	var (
		jsFields = newJSFields(millis)
		ticker   = time.NewTicker(jsFields.duration)
	)

	wg.Add(1)
	go func() {
		for {
			select {
			case <-jsFields.stopCh:
				ticker.Stop()
				wg.Done()
				return
			case <-ticker.C:
				f()
			}
		}
	}()
	return jsFields.stopExecute
}

// SetTimeout called function 'f' after 'millis' milliseconds after the call. It also returns
// a function which prevents execution SetTimeout
func SetTimeout(f func(), millis uint32, wg *sync.WaitGroup) func() {
	var (
		jsFields = newJSFields(millis)
		timer    = time.NewTimer(jsFields.duration)
	)

	wg.Add(1)
	go func() {
		for {
			select {
			case <-jsFields.stopCh:
				timer.Stop()
				wg.Done()
				return
			case <-timer.C:
				f()
				jsFields.closeChannel()
				wg.Done()
				return
			}
		}
	}()
	return jsFields.stopExecute
}
