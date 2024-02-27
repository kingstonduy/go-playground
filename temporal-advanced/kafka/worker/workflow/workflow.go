package workflow

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.temporal.io/sdk/workflow"
)

/**
 * The sample demonstrates how to deal with multiple signals that can come out of order and require actions
 * if a certain signal not received in a specified time interval.
 *
 * This specific sample receives three signals: Signal1, Signal2, Signal3. They have to be processed in the
 * sequential order, but they can be received out of order.
 * There are two timeouts to enforce.
 * The first one is the maximum time between signals.
 * The second limits the total time since the first signal received.
 *
 * A naive implementation of such use case would use a single loop that contains a Selector to listen on three
 * signals and a timer. Something like:

 *	for {
 *		selector := workflow.NewSelector(ctx)
 *		selector.AddReceive(workflow.GetSignalChannel(ctx, "Signal1"), func(c workflow.ReceiveChannel, more bool) {
 *			// Process signal1
 *		})
 *		selector.AddReceive(workflow.GetSignalChannel(ctx, "Signal2"), func(c workflow.ReceiveChannel, more bool) {
 *			// Process signal2
 *		}
 *		selector.AddReceive(workflow.GetSignalChannel(ctx, "Signal3"), func(c workflow.ReceiveChannel, more bool) {
 *			// Process signal3
 *		}
 *		cCtx, cancel := workflow.WithCancel(ctx)
 *		timer := workflow.NewTimer(cCtx, timeToNextSignal)
 *		selector.AddFuture(timer, func(f workflow.Future) {
 *			// Process timeout
 *		})
 * 		selector.Select(ctx)
 *	    cancel()
 *      // break out of the loop on certain condition
 *	}
 *
 *  The above implementation works. But it quickly becomes pretty convoluted if the number of signals
 *  and rules around order of their arrivals and timeouts increases.
 *
 *  The following example demonstrates an alternative approach. It receives signals in a separate goroutine.
 *  Each signal handler just updates a correspondent shared variable with the signal data.
 *  The main workflow function awaits the next step using `workflow.AwaitWithTimeout` using condition composed of
 *  the shared variables. This makes the main workflow method free from signal callbacks and makes the business logic
 *  clear.
 */

// SignalToSignalTimeout is them maximum time between signals
var SignalToSignalTimeout = 30 * time.Second

// FromFirstSignalTimeout is the maximum time to receive all signals
var FromFirstSignalTimeout = 60 * time.Second

type AwaitSignals struct {
	FirstSignalTime time.Time
	SignalReceived  bool
}

// Listen to signals Signal1, Signal2, and Signal3
func (a *AwaitSignals) Listen(ctx workflow.Context) {
	log := workflow.GetLogger(ctx)
	for {
		selector := workflow.NewSelector(ctx)
		selector.AddReceive(workflow.GetSignalChannel(ctx, workflow.GetInfo(ctx).WorkflowExecution.ID), func(c workflow.ReceiveChannel, more bool) {
			c.Receive(ctx, nil)
			a.SignalReceived = true
			log.Info("Signal Received")
		})

		selector.Select(ctx)
		if a.FirstSignalTime.IsZero() {
			a.FirstSignalTime = workflow.Now(ctx)
		}
	}
}

// AwaitSignalsWorkflow workflow definition
func AwaitSignalsWorkflow(ctx workflow.Context, transactionInfo string) (err error) {

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx = workflow.WithActivityOptions(ctx, ao)
	var a AwaitSignals

	workflow.ExecuteActivity(ctx, Produce, "localhost:9092", "topic1", "Hello Kafka")

	// Listen to signals in a different goroutine
	workflow.Go(ctx, a.Listen)

	// Wait for Signal
	err = workflow.Await(ctx, func() bool {
		return a.SignalReceived
	})
	// Cancellation
	if err != nil {
		return
	}

	return nil
}

func BoringActivity(ctx context.Context) error {
	timeout := rand.Intn(4)
	time.Sleep(time.Second * time.Duration(timeout))
	fmt.Printf("This Activity takes %d seconds to finish", timeout)

	return nil
}

func FunActivity(ctx context.Context) error {
	timeout := rand.Intn(4)
	time.Sleep(time.Second * time.Duration(timeout))
	fmt.Printf("Success💡💡💡💡💡💡💡", timeout)

	return nil
}

func Produce(ctx context.Context, bootstrapServer string, topic string, input string) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": bootstrapServer})
	if err != nil {
		panic(err)
	}

	// // Delivery report handler for produced messages
	// go func() {
	// 	for e := range p.Events() {
	// 		switch ev := e.(type) {
	// 		case *kafka.Message:
	// 			if ev.TopicPartition.Error != nil {
	// 				fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
	// 			} else {
	// 				fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
	// 			}
	// 		}
	// 	}
	// }()

	// Produce messages to topic (asynchronously)
	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          []byte(input),
	}, nil)

	// Wait for message deliveries
	p.Flush(15 * 1000)
	p.Close()
}
