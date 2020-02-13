package jstime

import (
	"sync"
	"time"
)

// SetInterval called function 'f' every 'millis' milliseconds. It also returns a channel
// that stops the function call 'f' by function ClearTimer(clear chan struct{})
func SetInterval(f func(), millis uint32, wg *sync.WaitGroup) chan struct{} {
	duration := time.Duration(millis) * time.Millisecond
	ticker := time.NewTicker(duration)
	clear := make(chan struct{})
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
	return clear
}

// SetTimeout called function 'f' after 'millis' milliseconds after the call. It also returns
// a channel which prevents execution SetTimeout by function ClearTimer(clear chan struct{})
func SetTimeout(f func(), millis uint32, wg *sync.WaitGroup) chan struct{} {
	duration := time.Duration(millis) * time.Millisecond
	timer := time.NewTimer(duration)
	clear := make(chan struct{})
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
				wg.Done()
			}
		}
	}()
	return clear
}

// ClearTimer stops execution of functions SetInterval() or SetTimeout()
func ClearTimer(clear chan struct{}) {
	clear <- struct{}{}
}
