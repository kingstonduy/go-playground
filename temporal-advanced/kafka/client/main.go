package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-gonic/gin"
	"github.com/pborman/uuid"
	"go.temporal.io/sdk/client"
)

type WorkflowInfo struct {
	ID    string `json:"ID`
	RunID string `json:"RunID`
}

// request, response to client. starts a workflow then wait for the workflow to finish
func Produce(c client.Client) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {
		options := client.StartWorkflowOptions{
			ID:        "my_workflow_id" + "-" + uuid.New(),
			TaskQueue: "taskqueue",
		}

		we, err := c.ExecuteWorkflow(context.Background(), options, "AwaitSignalsWorkflow")
		if err != nil {
			ctx.JSON(500, error.Error(err))
			return
		}

		err = we.Get(ctx, nil)
		if err != nil {
			ctx.JSON(500, error.Error(err))
			return
		}

		ctx.JSON(200, "Success")
	}
	return fn
}

func ParseWorkflowInfo(jsonStr string, info *WorkflowInfo) error {
	err := json.Unmarshal([]byte(jsonStr), &info)
	return err
}

// consume from kafka then signal the workflow to continue
func Consume(cl client.Client, bootstrapServer string, topic string) {

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

	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		router := gin.Default()

		router.POST("produce", Produce(c))

		router.Run("localhost:3000")
		defer wg.Done()
	}()

	wg.Add(1)
	go func() {
		Consume(c, "localhost:9092", "producer_topic")
		defer wg.Done()
	}()
	wg.Wait()

}
