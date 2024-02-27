package main

import (
	"fmt"
	"go-playground/temporal-advanced/async-http/worker/workflow"
	"log"
	"sync"

	"github.com/pborman/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// The client and worker are heavyweight objects that should be created once per process.
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	wg := sync.WaitGroup{}

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(i int) {
			w := worker.New(c, "taskqueue", worker.Options{
				BuildID:  uuid.New(),
				Identity: "worker" + fmt.Sprintf("%d", i),
			})

			w.RegisterWorkflow(workflow.AwaitSignalsWorkflow)
			w.RegisterActivity(workflow.BoringActivity)
			w.RegisterActivity(workflow.FunActivity)
			err = w.Run(worker.InterruptCh())
			if err != nil {
				log.Fatalln("Unable to start worker", err)
			}
			defer wg.Done()
		}(i)
	}
	wg.Wait()
}
