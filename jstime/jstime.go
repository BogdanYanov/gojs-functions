package jstime

import (
	"sync"
	"time"
)

type timeJSHelper struct {
	stopCh   chan struct{}
	mu       *sync.Mutex
	isClosed bool
}

func newTimeJSHelper() *timeJSHelper {
	return &timeJSHelper{
		stopCh:   make(chan struct{}),
		mu:       &sync.Mutex{},
		isClosed: false,
	}
}

func (helper *timeJSHelper) isChannelClosed() bool {
	helper.mu.Lock()
	defer helper.mu.Unlock()
	return helper.isClosed
}

func (helper *timeJSHelper) stopExecute() {
	if !helper.isChannelClosed() {
		helper.stopCh <- struct{}{}
		helper.closeChannel()
	}
}

func (helper *timeJSHelper) closeChannel() {
	helper.mu.Lock()
	close(helper.stopCh)
	helper.isClosed = true
	helper.mu.Unlock()
}

// SetInterval called function 'f' every 'millis' milliseconds. It also returns a function
// that stops the function call 'f'
func SetInterval(f func(), duration time.Duration) func() {
	var (
		timeJSHelper = newTimeJSHelper()
		ticker             = time.NewTicker(duration)
	)

	go func() {
		for {
			select {
			case <-timeJSHelper.stopCh:
				ticker.Stop()
				return
			case <-ticker.C:
				f()
			}
		}
	}()
	return timeJSHelper.stopExecute
}

// SetTimeout called function 'f' after 'millis' milliseconds after the call. It also returns
// a function which prevents execution SetTimeout
func SetTimeout(f func(), duration time.Duration) func() {
	var (
		timeJSHelper = newTimeJSHelper()
		timer              = time.NewTimer(duration)
	)

	go func() {
		for {
			select {
			case <-timeJSHelper.stopCh:
				timer.Stop()
				return
			case <-timer.C:
				f()
				timeJSHelper.closeChannel()
				return
			}
		}
	}()
	return timeJSHelper.stopExecute
}
