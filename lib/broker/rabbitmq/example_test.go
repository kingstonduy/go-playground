package main

import (
	"fmt"
	"log"
	"math/rand"
	"sync"
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
)

type resquest struct {
	Message string
}

type response struct {
	Message string
}

func Test(t *testing.T) {
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
