package main

import (
	"log"
	"time"
)

func main() {
	for {
		rand := newRand()

		// start twin reader and writer goroutines
		readerDone := make(chan bool)
		go reader(readerDone, rand)
		writerDone := make(chan bool)
		go writer(writerDone, rand)

		// block until one of reader or writer signals done
		select {
		case <-readerDone:
			// reader is done, close the writer
			writerDone <- true
		case <-writerDone:
			// writer is done, close the reader
			readerDone <- true
		}

		// close both channels
		close(readerDone)
		close(writerDone)
	}
}

func reader(readerDone chan bool, rand *randWrapper) {
READ_LOOP:
	for {
		select {
		case <-readerDone:
			break READ_LOOP
		default:
			log.Println("Reader doing work")
			sleepTime := rand.sleepTime()
			log.Printf("Reader sleeping for %d seconds", sleepTime/time.Second)
			time.Sleep(sleepTime)
		}
	}
}

func writer(writerDone chan bool, rand *randWrapper) {
WRITE_LOOP:
	for {
		select {
		case <-writerDone:
			break WRITE_LOOP
		default:
			log.Println("Writer doing work")
			sleepTime := rand.sleepTime()
			log.Printf("Writer sleeping for %d seconds", sleepTime/time.Second)
			time.Sleep(sleepTime)
		}
	}
}
