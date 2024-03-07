package main

import (
	"context"
	"encoding/json"
	"sync"

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
