package main

import (
	"github.com/BogdanYanov/gojs-functions/jstime"
	"log"
	"sync"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	var (
		wg       = &sync.WaitGroup{}
		interval = 500 * time.Millisecond
		timeout1 = 5 * time.Second
		timeout2 = 4500 * time.Millisecond
		timeout3 = 4100 * time.Millisecond
	)

	wg.Add(1) // first - ADD
	// Every 500 milliseconds display a message
	stopI := jstime.SetInterval(func() {
		log.Println("Tick")
	}, interval)

	wg.Add(1) // second - ADD
	// Stop message output after 5 seconds
	stopT := jstime.SetTimeout(func() {
		log.Println("Stop interval after 5 second")
		stopI()
		wg.Done() // second - DONE
	}, timeout1)

	time.Sleep(4 * time.Second)
	log.Println("Stop interval after 4 second")
	stopI()
	wg.Done() // first - DONE

	wg.Add(1) // third - ADD
	stopT = jstime.SetTimeout(func() {
		log.Println("Execute new timeout after 8,5 second")
	}, timeout2)

	wg.Add(1) // fourth - ADD
	_ = jstime.SetTimeout(func() {
		log.Println("Stop new timeout after 8,1 second")
		stopT()
		wg.Done() // fourth - DONE
	}, timeout3)

	// Don`t stop message output after 4 seconds
	log.Println("Stop new timeout after 4 second")
	stopT()
	wg.Done() // third - DONE
	stopT()
	stopT()

	time.Sleep(1500 * time.Millisecond)

	// Stop message output
	log.Println("Stop interval after 5,5 second")
	stopI()
	stopI()
	stopI()

	wg.Wait()
}
