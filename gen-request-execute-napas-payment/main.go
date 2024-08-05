package main

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/IBM/sarama"
	"github.com/google/uuid"
)

type Command struct {
	CommandID            string `json:"COMMAND_ID"`
	AggregateID          string `json:"AGGREGATE_ID"`
	CommandType          string `json:"COMMAND_TYPE"`
	AggregateType        string `json:"AGGREGATE_TYPE"`
	Payload              string `json:"PAYLOAD"`
	Processed            int8   `json:"PROCESSED"`
	ProcessAt            int64  `json:"PROCESS_AT"`
	ReplyTo              string `json:"REPLY_TO"`
	TransactionCreatedAt int64  `json:"TRANSACTION_CREATED_AT"`
}

type ExecuteNapasPaymentRequest struct {
	ClientTransId           string                        `json:"clientTransId" validate:"required"`
	PaymentMethod           string                        `json:"paymentMethod" validate:"required"`
	PaymentChannel          string                        `json:"paymentChannel" validate:"required"`
	Channel                 string                        `json:"channel" validate:"required"`
	BranchCode              string                        `json:"branchCode"`
	FromAccountNumber       string                        `json:"fromAccountNumber"`
	ToBankCode              string                        `json:"toBankCode"`
	ToAccountNumber         string                        `json:"toAccountNumber"`
	ToCardNumber            string                        `json:"toCardNumber"`
	ToCreditCard            ExecuteNapasPaymentCreditCard `json:"toCreditCard"`
	Amount                  int64                         `json:"amount" validate:"required"`
	Currency                string                        `json:"currency" validate:"required"`
	Remark                  string                        `json:"remark" validate:"required"`
	BenefitCustomerName     string                        `json:"benefitCustomerName" validate:"required"`
	NapasRefNumber          string                        `json:"napasRefNumber" validate:"required"`
	CustomerID              string                        `json:"customerId" validate:"required"`
	MerchantCategoryCode    string                        `json:"merchantCategoryCode"`
	AcceptorNameAndLocation string                        `json:"acceptorNameAndLocation"`
}

type ExecuteNapasPaymentCreditCard struct {
	EncryptedCardNo string `json:"encryptedCardNo"`
	EncryptedKey    string `json:"encryptedKey"`
	CardMasking     string `json:"cardMasking"`
}

type Trace struct {
	From     string `json:"frm" xml:"frm"`
	To       string `json:"to" xml:"to"`
	Cid      string `json:"cid" xml:"cid" validate:"required"`
	Sid      string `json:"sid" xml:"sid"`
	Cts      int64  `json:"cts" xml:"cts" validate:"required"`
	Sts      int64  `json:"sts" xml:"sts"`
	Dur      int64  `json:"dur" xml:"dur"`
	Username string `json:"userName" xml:"userName"`
	// MessageType string   `json:"messageType" xml:"messageType"`
	// ReplyTo     []string `json:"replyTo" xml:"replyTo"`
}

type Request[T any] struct {
	Trace Trace `json:"trace" xml:"trace" validate:"required"`
	Data  T     `json:"data" xml:"data"`
}

func (r *Request[T]) SetTrace(trace Trace) {
	r.Trace = trace
}

func (r *Request[T]) GetTrace() Trace {
	return r.Trace
}

type Result struct {
	StatusCode int         `json:"statusCode" xml:"statusCode"`
	Code       string      `json:"code" xml:"code"`
	Message    string      `json:"message" xml:"message"`
	Details    interface{} `json:"details" xml:"details"`
}

type Response[T any] struct {
	Result Result `json:"result" xml:"result"`
	Trace  Trace  `json:"trace" xml:"trace"`
	Data   T      `json:"data" xml:"data"`
}

func ProduceStruct() {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer([]string{"10.96.24.141:9093"}, config)
	if err != nil {
		panic(err)
	}

	// Trap SIGINT to trigger a graceful shutdown.
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	var (
		wg                                  sync.WaitGroup
		enqueued, successes, producerErrors int
	)

	wg.Add(1)
	go func() {
		defer wg.Done()
		for range producer.Successes() {
			successes++
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range producer.Errors() {
			log.Println(err)
			producerErrors++
		}
	}()

	done := make(chan bool)
	go func() {
		<-signals
		fmt.Println("Received interrupt signal. Shutting down...")
		producer.AsyncClose() // Trigger a shutdown of the producer.
		close(done)
	}()

	numLoop := 1
ProducerLoop:
	for i := 1; i <= numLoop; i++ {

		rand_uuid := uuid.New().String()
		var event Command
		event.CommandType = "INITIALIZE_TRANSACTION"
		event.AggregateID = rand_uuid
		event.ReplyTo = "SAGA.NAPAS_FAST_FUND.RES.V1"

		payload := Request[ExecuteNapasPaymentRequest]{
			Trace: Trace{
				From: "fast-fund-service",
				To:   "saga",
				Cid:  uuid.New().String(),
				Sid:  uuid.New().String(),
				Cts:  time.Now().UnixMilli(),
			},
			Data: ExecuteNapasPaymentRequest{
				ClientTransId:  randomString(),
				PaymentMethod:  "ACCOUNT",
				PaymentChannel: "IB",
				Channel:        "OPENAPI",
				// BranchCode:              "",
				FromAccountNumber: "0037100023298001",
				ToBankCode:        "970406",
				ToAccountNumber:   "0129837294",
				// ToCardNumber:            "",
				// ToCreditCard:            "",
				Amount:              10000,
				Currency:            "VND",
				Remark:              "chuyentiennapas",
				BenefitCustomerName: "NGUYENVANTEST",
				NapasRefNumber:      "4129OCBB8804306106",
				CustomerID:          "7693058",
				// MerchantCategoryCode:    "",
				// AcceptorNameAndLocation: "",
			},
		}
		payloadBytes, _ := json.Marshal(payload)
		event.Payload = string(payloadBytes)
		event.CommandID = uuid.New().String()

		msg, _ := json.Marshal(event)
		log.Println(rand_uuid)
		message := &sarama.ProducerMessage{Topic: "NAPAS_FAST_FUND.SAGA.REQ.V1", Value: sarama.StringEncoder(msg)}
		select {
		case producer.Input() <- message:
			enqueued++
		case <-done:
			break ProducerLoop
		}
	}
	producer.AsyncClose() // Trigger a shutdown of the producer.
	// Wait for all go routines to finish
	wg.Wait()

	log.Printf("Successfully produced: %d; errors: %d\n", successes, producerErrors)
}

func randomString() string {
	length := 8
	if length <= 0 {
		return ""
	}
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	return base64.URLEncoding.EncodeToString(bytes)[:length]
}
