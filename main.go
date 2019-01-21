package main

import (
	"log"
	"sync"
	"time"
)

func main() {
	rand := newRand()
	for i := 0; true; i++ {
		log.Printf("Starting twin goroutines")
		done := newDone()

		go workLoop(read, done, rand)
		go workLoop(write, done, rand)

		// block until we're done
		select {
		case <-done.done:
			log.Printf("Someone signalled done")
		}

		time.Sleep(10 * time.Second)
	}
}

func read() {
	log.Println("Reader doing work")
}

func write() {
	log.Println("Writer doing work")
}

func workLoop(doWork func(), done *done, rand *randWrapper) {
WORK_LOOP:
	for {
		select {
		case <-done.done:
			break WORK_LOOP
		default:
			doWork()

			snakeEyes := rand.snakeEyes()
			if snakeEyes {
				log.Printf("Snake eyes! Signalling done.")
				done.signalDone()
			}

			sleepTime := rand.sleepTime()
			log.Printf("Sleeping for %d seconds", sleepTime/time.Second)
			time.Sleep(sleepTime)
		}
	}
}

type done struct {
	done chan struct{}
	once sync.Once
}

func newDone() *done {
	ch := make(chan struct{})
	return &done{done: ch}
}

func (d *done) signalDone() {
	// you can only close a channel once or Go panics
	d.once.Do(func() {
		close(d.done)
	})
}
