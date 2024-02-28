package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type resquest struct {
	Message string
}

type response struct {
	Message string
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func fib(n resquest) response {
	// sleep for random 0-5 second
	time.Sleep(time.Duration(rand.Intn(20)))
	return response{Message: "Hello, " + n.Message}
}

func ConsumeAndPublish(topic string, url string) {
	conn, err := amqp.Dial(url)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		topic, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	concurrency := 10000
	err = ch.Qos(
		concurrency, // prefetch count
		0,           // prefetch size
		false,       // global
	)
	failOnError(err, "Failed to set QoS")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	failOnError(err, "Failed to register a consumer")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var wg sync.WaitGroup // used to coordinate when they are done, ie: if rabbit conn was closed
	wg.Add(concurrency)

	for x := 0; x < concurrency; x++ {
		go func() {
			defer wg.Done()
			for d := range msgs {
				log.Printf(" [*] Awaiting RPC requests")

				var req resquest
				var res response

				_ = json.Unmarshal(d.Body, &req)

				log.Printf(" [.] request %+v\n", req)

				res = fib(req)
				resJSON, _ := json.Marshal(res)

				err = ch.PublishWithContext(ctx,
					"",        // exchange
					d.ReplyTo, // routing key
					false,     // mandatory
					false,     // immediate
					amqp.Publishing{
						ContentType:   "text/plain",
						CorrelationId: d.CorrelationId,
						Body:          []byte(resJSON),
					})
				failOnError(err, "Failed to publish a message")

				d.Ack(false)
			}
		}()
	}
	wg.Wait() // when all goroutine's exit, the app exits
}

func main() {
	ConsumeAndPublish("rpc_queue11111", "amqp://guest:guest@localhost:5673/")
}
