package main

import (
	"math/rand"
	"sync"
	"time"
)

const (
	maxSleepTime = 5
)

type randWrapper struct {
	// https://golang.org/pkg/math/rand/#NewSource
	// "this source is not safe for concurrent use by multiple goroutines"
	// so we use a mutex here
	rand *rand.Rand
	mux  sync.Mutex
}

func newRand() *randWrapper {
	time := time.Now().UnixNano()
	source := rand.NewSource(time)
	rand := rand.New(source)
	return &randWrapper{rand: rand}
}

func (rw *randWrapper) sleepTime() time.Duration {
	rw.mux.Lock()
	defer rw.mux.Unlock()
	return time.Duration(rw.rand.Intn(maxSleepTime)) * time.Second
}
