package jstime

import (
	"sync"
	"testing"
	"time"
)

func TestSetInterval(t *testing.T) {
	var (
		timer1  *time.Timer
		timer2  *time.Timer
		counter int
		wg      *sync.WaitGroup
	)

	counter = 0
	wg = &sync.WaitGroup{}
	timer1 = time.NewTimer(2000 * time.Millisecond)
	timer2 = time.NewTimer(2500 * time.Millisecond)

	stop := SetInterval(func() {
		counter++
	}, 250, wg)

OuterLoop:
	for {
		select {
		case <-timer2.C:
			if counter != 8 {
				t.Errorf("Counter = %d, want - %d", counter, 8)
				stop()
			}
			break OuterLoop
		case <-timer1.C:
			stop()
		default:
		}
	}
	wg.Wait()
}

func TestSetTimeout(t *testing.T) {
	var (
		timer1  *time.Timer
		timer2  *time.Timer
		timer3  *time.Timer
		counter int
		wg      *sync.WaitGroup
	)

	counter = 0
	wg = &sync.WaitGroup{}
	timer1 = time.NewTimer(1000 * time.Millisecond)
	timer2 = time.NewTimer(2500 * time.Millisecond)
	timer3 = time.NewTimer(3000 * time.Millisecond)

	stopI := SetInterval(func() {
		counter++
	}, 250, wg)

	stopT := SetTimeout(func() {
		stopI()
	}, 2000, wg)

	_ = SetTimeout(func() {
		stopT()
	}, 1500, wg)

OuterLoop:
	for {
		select {
		case <-timer1.C:
			stopT()
		case <-timer2.C:
			stopI()
		case <-timer3.C:
			if counter != 10 {
				t.Errorf("Counter = %d, want - %d", counter, 8)
				stopI()
			}
			break OuterLoop
		default:
		}
	}
	wg.Wait()
}
