package main

import (
	"context"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/pborman/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

func RequestAndReply(n int) (res string, err error) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5673/")
	if err != nil {
		log.Panicf("%s: Failed to connect to RabbitMQ", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("%s: Failed to open a channel", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		log.Panicf("%s: Failed to declare a queue", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Panicf("%s: Failed to register a consumer", err)
	}

	corrId := uuid.New()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",          // exchange
		"rpc_queue", // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(strconv.Itoa(n)),
		})
	if err != nil {
		log.Panicf("%s: Failed to publish a message", err)
	}

	for d := range msgs {
		if corrId == d.CorrelationId {
			res = string(d.Body)
			if err != nil {
				log.Panicf("%s: Failed to publish a message", err)
			}
			break
		}
	}

	return
}

func main() {
	wg := sync.WaitGroup{}

	for i := 1; i <= 1000; i++ {
		wg.Add(1)
		go func() {
			n := rand.Intn(20)
			log.Printf(" [x] Requesting fib(%d)", n)
			res, err := RequestAndReply(n)
			if err != nil {
				log.Panicf("%s: Failed to handle RPC request", err)
			}

			log.Printf(" [.] Response for n=%d:%s", n, res)
			defer wg.Done()
		}()
	}

	wg.Wait()
}
