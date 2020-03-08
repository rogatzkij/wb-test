package main

import (
	"sync"
)

func main() {
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go Major(wg)
	wg.Wait()
}
