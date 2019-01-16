package main

import (
	"log"
	"time"
)

type doneSignal struct{}

func main() {
	rand := newRand()
	for i := 0; true; i++ {
		log.Printf("Starting twin goroutines")
		done := make(chan doneSignal)
		go reader(done, rand)
		go writer(done, rand)

		// block until we're done
		select {
		case <-done:
			log.Printf("Someone signalled done, closing channel")
			close(done)
		}

		time.Sleep(10 * time.Second)
	}
}

func reader(done chan doneSignal, rand *randWrapper) {
READ_LOOP:
	for {
		select {
		case <-done:
			break READ_LOOP
		default:
			log.Println("Reader doing work")

			snakeEyes := rand.snakeEyes()
			if snakeEyes {
				log.Printf("Snake eyes! Signalling reader done.")
				done <- doneSignal{}
			}

			sleepTime := rand.sleepTime()
			log.Printf("Reader sleeping for %d seconds", sleepTime/time.Second)
			time.Sleep(sleepTime)
		}
	}
}

func writer(done chan doneSignal, rand *randWrapper) {
WRITE_LOOP:
	for {
		select {
		case <-done:
			break WRITE_LOOP
		default:
			log.Println("Writer doing work")

			snakeEyes := rand.snakeEyes()
			if snakeEyes {
				log.Printf("Snake eyes! Signalling writer done.")
				done <- doneSignal{}
			}

			sleepTime := rand.sleepTime()
			log.Printf("Writer sleeping for %d seconds", sleepTime/time.Second)
			time.Sleep(sleepTime)
		}
	}
}
