package main

import (
	"log"
	"time"
)

func main() {
	rand := newRand()
	for {
		log.Printf("Starting twin goroutines")

		// start twin reader and writer goroutines
		readerDone := make(chan bool)
		go reader(readerDone, rand)
		writerDone := make(chan bool)
		go writer(writerDone, rand)

		// block until one of reader or writer signals done
		select {
		case <-readerDone:
			log.Printf("Reader signalled done, closing writer")
			writerDone <- true
		case <-writerDone:
			log.Printf("Writer signalled done, closing writer")
			readerDone <- true
		}

		// close both channels
		// BUG: it's entirely possible that the reader or writer is still going
		// 2019/01/16 22:29:17 Reader signalled done, closing writer
		// 2019/01/16 22:29:17 Reader doing work
		// 2019/01/16 22:29:17 Reader doing work
		// 2019/01/16 22:29:17 Reader doing work
		// 2019/01/16 22:29:17 Reader doing work
		// 2019/01/16 22:29:17 Reader doing work
		// 2019/01/16 22:29:17 Snake eyes! Signalling reader done.
		// 2019/01/16 22:29:17 Writer doing work
		// panic: send on closed channel
		close(readerDone)
		close(writerDone)

		time.Sleep(10 * time.Second)
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

			snakeEyes := rand.snakeEyes()
			if snakeEyes {
				log.Printf("Snake eyes! Signalling reader done.")
				readerDone <- true
			}

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

			snakeEyes := rand.snakeEyes()
			if snakeEyes {
				log.Printf("Snake eyes! Signalling writer done.")
				writerDone <- true
			}

			sleepTime := rand.sleepTime()
			log.Printf("Writer sleeping for %d seconds", sleepTime/time.Second)
			time.Sleep(sleepTime)
		}
	}
}
