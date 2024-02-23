package main

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
	"sync"
)

type Subscriber chan interface{}

type TopicFilter func(interface{}) bool

type Publisher struct {
	m           sync.Mutex
	buffer      int
	Subscribers map[Subscriber]TopicFilter
}

func NewPublisher(buffer int) *Publisher {
	return &Publisher{
		buffer:      buffer,
		Subscribers: make(map[Subscriber]TopicFilter),
	}
}

func (p *Publisher) Publish(msg string) {
	wg := sync.WaitGroup{}
	for sub, topic := range p.Subscribers {
		if !topic(msg) {
			continue
		}
		wg.Add(1)
		go func() {
			sub <- msg
			defer wg.Done()
		}()
	}
	wg.Wait()
}

// Subscribe all topics
func (p *Publisher) Subscribe(topic TopicFilter) chan interface{} {
	ch := make(chan interface{}, p.buffer)
	p.m.Lock()
	p.Subscribers[ch] = topic
	p.m.Unlock()
	return ch
}

// Unsubscribe
func (p *Publisher) Evict(sub chan interface{}) {
	p.m.Lock()
	defer p.m.Unlock()
	close(sub)
}

func (p *Publisher) Close() {
	p.m.Lock()
	defer p.m.Unlock()
	for sub := range p.Subscribers {
		delete(p.Subscribers, sub)
		close(sub)
	}
}

func main() {
	p := NewPublisher(1000)
	defer p.Close()

	all := p.Subscribe(func(v interface{}) bool {
		return true
	})

	group1 := p.Subscribe(func(v interface{}) bool {
		if s, ok := v.(string); ok {
			return strings.Contains(s, "1")
		}
		return false
	})

	group2 := p.Subscribe(func(v interface{}) bool {
		if s, ok := v.(string); ok {
			return strings.Contains(s, "2")
		}
		return false
	})

	group3 := p.Subscribe(func(v interface{}) bool {
		if s, ok := v.(string); ok {
			return strings.Contains(s, "3")
		}
		return false
	})

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		for {
			var input string
			fmt.Scan(&input)
			if input == "exit" {
				break
			}
			nBig, err := rand.Int(rand.Reader, big.NewInt(3))
			if err != nil {
				panic(err)
			}
			n := nBig.Int64() + 1

			input = fmt.Sprintf("[GROUP %d]: %s", n, input)
			fmt.Println("Input:" + input)
			p.Publish(input)
		}
		defer wg.Done()
	}()

	wg.Add(1)
	go func() {
		for _ = range all {
			fmt.Printf("Received all: %v\n", <-all)
		}
		defer wg.Done()
	}()

	wg.Add(1)
	go func() {
		for _ = range group1 {
			fmt.Printf("Received group 1: %v\n", <-group1)
		}
		defer wg.Done()
	}()

	wg.Add(1)
	go func() {
		for _ = range group2 {
			fmt.Printf("Received group 2: %v\n", <-group2)
		}
		defer wg.Done()
	}()

	wg.Add(1)
	go func() {
		for _ = range group3 {
			fmt.Printf("Received group 3: %v\n", <-group3)
		}
		defer wg.Done()
	}()

	wg.Wait()

}
