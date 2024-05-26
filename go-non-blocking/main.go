package main

import (
	"fmt"
	"sync"
)

func doSth(id string) {
	for i := 1; i <= 100; i++ {
		fmt.Printf("Go routine=%s. Epoch=%d\n", id, i)
	}
}

func main() {
	numGoroutines := 100
	wg := sync.WaitGroup{}

	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			doSth(id)
		}(fmt.Sprintf("%d", i))
	}

	wg.Wait()
}
