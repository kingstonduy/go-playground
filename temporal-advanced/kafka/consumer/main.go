package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.temporal.io/sdk/client"
)

type WorkflowInfo struct {
	ID           string `json:"ID`
	RunID        string `json:"RunID`
	ActivityName string `json:ActivityName`
}

// consume from kafka then signal the workflow to continue
func Consume[T any, K any](request T, response *K, cl client.Client, bootstrapServer string, topic string) {

	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": bootstrapServer,
		"group.id":          "myGroup",
		"auto.offset.reset": "latest",
	})

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{topic}, nil)

	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			var workflowInfo WorkflowInfo
			err := json.Unmarshal([]byte(string(msg.Value)), &workflowInfo)
			if err != nil {
				log.Println(err)
				return
			}

			signalName := workflowInfo.ID
			err = cl.SignalWorkflow(context.Background(), workflowInfo.ID, workflowInfo.RunID, signalName, nil)
			if err != nil {
				log.Println(err)
				return
			}

		} else {
			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
			break
		}
	}

	c.Close()
}

func main() {
	fmt.Println("Hello world")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("quit (%v)\n", <-sig)
}
