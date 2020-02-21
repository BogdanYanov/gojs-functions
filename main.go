package main

import (
	"github.com/BogdanYanov/gojs-functions/jstime"
	"log"
	"sync"
	"time"
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	wg := &sync.WaitGroup{}

	// Every 500 milliseconds display a message
	stopI := jstime.SetInterval(func() {
		log.Println("Tick")
	}, 500, wg)

	// Stop message output after 5 seconds
	stopT := jstime.SetTimeout(func() {
		log.Println("Stop interval after 5 second")
		stopI()
	}, 5000, wg)

	time.Sleep(4000 * time.Millisecond)
	log.Println("Stop interval after 4 second")
	stopI()

	stopT = jstime.SetTimeout(func() {
		log.Println("Test new timeout")
	}, 4500, wg)

	_ = jstime.SetTimeout(func() {
		log.Println("Stop new timeout after 4100")
		stopT()
	}, 4100, wg)

	// Don`t stop message output after 4 seconds
	log.Println("Stop new timeout after 4 second")
	stopT()
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
