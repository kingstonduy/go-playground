package main

import (
	"context"
	"log"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/pborman/uuid"
	"go.temporal.io/sdk/client"
)

type MessageQueue struct {
	index map[string]chan string
	m     sync.Mutex
}

type WorkflowInfo struct {
	ID    string `json:"ID`
	RunID string `json:"RunID`
}

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

func Consume(c client.Client) gin.HandlerFunc {
	fn := func(ctx *gin.Context) {

		var workflowInfo WorkflowInfo
		err := ctx.BindJSON(&workflowInfo)
		if err != nil {
			log.Println(err)
			return
		}

		signalName := workflowInfo.ID
		err = c.SignalWorkflow(context.Background(), workflowInfo.ID, workflowInfo.RunID, signalName, nil)
		if err != nil {
			log.Println(err)
			ctx.JSON(500, err)
			return
		}

		ctx.JSON(200, "Success sending signal")
	}

	return fn
}

func main() {

	// The client is a heavyweight object that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		panic(err)
	}

	router := gin.Default()

	router.POST("produce", Produce(c))
	router.POST("comsume", Consume(c))

	router.Run("localhost:3000")
}
