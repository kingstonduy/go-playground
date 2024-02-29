package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/pborman/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type resquest struct {
	Message string
}

type response struct {
	Message string
}

func RequestAndReply[T any, K any](req T, res *K, topic string, conn *amqp.Connection) error {
	ch, err := conn.Channel()
	if err != nil {
		log.Panicf("%s: Failed to open a channel", err)
		return err
	}
	defer ch.Close()

	err = ch.Qos(
		100000, // prefetch count
		0,      // prefetch size
		false,  // global
	)

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
		return err
	}

	corrId := uuid.New()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	inputString, err := json.Marshal(req)
	if err != nil {
		log.Panicf("Failed to convert object to JSON: %s", err)
		return err
	}

	err = ch.PublishWithContext(ctx,
		"",    // exchange
		topic, // routing key
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(inputString),
		})
	if err != nil {
		log.Panicf("%s: Failed to publish a message", err)
		return err
	}

	for d := range msgs {
		if corrId == d.CorrelationId {
			err = json.Unmarshal(d.Body, res)
			if err != nil {
				log.Panicf("%s: Failed to convert json to  object", err)
				return err
			}
			break
		}
	}

	return nil
}

func main() {
	url := "amqp://guest:guest@localhost:5673/"
	conn, err := amqp.Dial(url)
	if err != nil {
		log.Panicf("%s: Failed to connect to RabbitMQ", err)
		return
	}
	defer conn.Close()

	wg := sync.WaitGroup{}
	for i := 1; i <= 2046; i++ {
		wg.Add(1)
		go func() {
			var request resquest
			var response response

			request.Message = fmt.Sprintf("%d", rand.Intn(100000))
			log.Printf(" [x] Request: %+v ", request)

			err := RequestAndReply(request, &response, "rpc_queue11111", conn)
			if err != nil {
				log.Panicf("%s: Failed to handle RPC request", err)
			}

			log.Printf(" [.] Response for n = %s:  %+v\n", request.Message, response)
			defer wg.Done()
		}()
	}

	wg.Wait()
}
