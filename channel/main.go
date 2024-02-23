package main

import (
	"fmt"
	"sync"
)

func main() {
	var links []int

	for i := 1; i <= 100; i++ {
		links = append(links, i)
	}

	linkChannel := make(chan int, 5) // Buffer for all links

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		for _, link := range links {
			for len(linkChannel) == cap(linkChannel) {
				<-linkChannel
			}
			linkChannel <- link
		}
		close(linkChannel) // Close the channel when done
		defer wg.Done()
	}()

	wg.Add(1)
	go func() {
		for link := range linkChannel { // Receive until closed
			fmt.Println(link)
		}
		defer wg.Done()
	}()

	wg.Wait()
}
