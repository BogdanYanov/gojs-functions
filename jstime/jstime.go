package jstime

import (
	"sync"
	"time"
)

// SetInterval called function 'f' every 'millis' milliseconds. It also returns a channel
// that stops the function call 'f' by function ClearTimer(clear chan struct{})
func SetInterval(f func(), millis uint32, wg *sync.WaitGroup) func() {
	var (
		duration time.Duration
		ticker *time.Ticker
		clear chan struct{}
		isClosed bool
	)

	duration = time.Duration(millis) * time.Millisecond
	ticker = time.NewTicker(duration)
	clear = make(chan struct{})
	wg.Add(1)
	go func() {
		for {
			select {
			case <-clear:
				ticker.Stop()
				close(clear)
				wg.Done()
				return
			case <-ticker.C:
				f()
			}
		}
	}()
	return func() {
		if !isClosed {
			clear <- struct{}{}
			isClosed = true
		}
	}
}

// SetTimeout called function 'f' after 'millis' milliseconds after the call. It also returns
// a channel which prevents execution SetTimeout by function ClearTimer(clear chan struct{})
func SetTimeout(f func(), millis uint32, wg *sync.WaitGroup) func() {
	var (
		duration time.Duration
		timer *time.Timer
		clear chan struct{}
		isClosed bool
	)

	duration = time.Duration(millis) * time.Millisecond
	timer = time.NewTimer(duration)
	clear = make(chan struct{})
	wg.Add(1)
	go func() {
		for {
			select {
			case <-clear:
				timer.Stop()
				close(clear)
				wg.Done()
				return
			case <-timer.C:
				f()
				close(clear)
				isClosed = true
				wg.Done()
				return
			}
		}
	}()
	return func() {
		if !isClosed {
			clear <- struct{}{}
			isClosed = true
		}
	}
}