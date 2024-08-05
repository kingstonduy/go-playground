package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/IBM/sarama"
)

func Produce() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer([]string{"10.96.24.141:9093"}, config)
	if err != nil {
		panic(err)
	}

	// Trap SIGINT to trigger a graceful shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var (
		wg                                  sync.WaitGroup
		enqueued, successes, producerErrors int
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range producer.Successes() {
			successes++
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range producer.Errors() {
			log.Println(err)
			producerErrors++
		}
	}()

	done := make(chan bool)
	go func() {
		<-signals
		fmt.Println("Received interrupt signal. Shutting down...")
		producer.AsyncClose() // Trigger a shutdown of the producer.
		close(done)
	}()

ProducerLoop:
	for i := 1; i <= 10; i++ {
		message := &sarama.ProducerMessage{Topic: "test-duydk3", Value: sarama.StringEncoder("testing 123")}
		select {
		case producer.Input() <- message:
			enqueued++
		case <-done:
			break ProducerLoop
		}
	}
	producer.AsyncClose() // Trigger a shutdown of the producer.
	// Wait for all go routines to finish
	wg.Wait()

	log.Printf("Successfully produced: %d; errors: %d\n", successes, producerErrors)
}
