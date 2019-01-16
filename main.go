package main

import (
	"fmt"
	"time"
)

func main() {
	for {
		// start twin reader and writer goroutines
		readerDone := make(chan bool)
		go reader(readerDone)
		writerDone := make(chan bool)
		go writer(writerDone)

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

func reader(readerDone chan bool) {
READ_LOOP:
	for {
		select {
		case <-readerDone:
			break READ_LOOP
		default:
			fmt.Println("Reader doing work")
			time.Sleep(1 * time.Second)
		}
	}
}

func writer(writerDone chan bool) {
WRITE_LOOP:
	for {
		select {
		case <-writerDone:
			break WRITE_LOOP
		default:
			fmt.Println("Writer doing work")
			time.Sleep(3 * time.Second)
		}
	}
}
