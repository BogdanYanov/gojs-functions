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
	clearI := jstime.SetInterval(func() {
		log.Println("Tick")
	}, 500, wg)

	// Stop message output after 5 seconds
	clearT := jstime.SetTimeout(func() {
		jstime.ClearTimer(clearI)
	}, 5000, wg)

	time.Sleep(4000 * time.Millisecond)

	// Don`t stop message output after 4 seconds
	jstime.ClearTimer(clearT)

	time.Sleep(1500 * time.Millisecond)

	// Stop message output
	jstime.ClearTimer(clearI)

	wg.Wait()
}
