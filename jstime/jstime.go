package jstime

import (
	"sync"
	"time"
)

// SetInterval called function 'f' every 'millis' milliseconds. It also returns a function
// that stops the function call 'f'
func SetInterval(f func(), millis uint32, wg *sync.WaitGroup) func() {
	var (
		duration time.Duration
		ticker   *time.Ticker
		clear    chan struct{}
		mu       *sync.Mutex
		isClosed bool
	)

	duration = time.Duration(millis) * time.Millisecond
	ticker = time.NewTicker(duration)
	clear = make(chan struct{})
	mu = &sync.Mutex{}

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
		mu.Lock()
		if !isClosed {
			clear <- struct{}{}
			isClosed = true
		}
		mu.Unlock()
	}
}

// SetTimeout called function 'f' after 'millis' milliseconds after the call. It also returns
// a function which prevents execution SetTimeout
func SetTimeout(f func(), millis uint32, wg *sync.WaitGroup) func() {
	var (
		duration time.Duration
		timer    *time.Timer
		clear    chan struct{}
		mu       *sync.Mutex
		isClosed bool
	)

	duration = time.Duration(millis) * time.Millisecond
	timer = time.NewTimer(duration)
	clear = make(chan struct{})
	mu = &sync.Mutex{}

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
				mu.Lock()
				isClosed = true
				mu.Unlock()
				wg.Done()
				return
			}
		}
	}()
	return func() {
		mu.Lock()
		if !isClosed {
			clear <- struct{}{}
			isClosed = true
		}
		mu.Unlock()
	}
}
