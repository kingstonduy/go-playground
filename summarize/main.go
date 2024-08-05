package main

type Command struct {
	AggregateID   string `json:"aggregateID"`
	AggregateType string `json:"aggregateType"`
	CommandID     string `json:"commandID"`
	CommandType   string `json:"commandType"`
	Payload       string `json:"payload"`
}

type BrokerHandler interface {
}

type BrokerServer interface {
	Publish(topic string, body []byte, key []byte, headers map[string][]string) error
	Subscribe(topic string, handler BrokerHandler) error
}

func NewKafkaServer() BrokerServer {
	return &KafkaServer{}
}

type KafkaServer struct{}

// Publish implements BrokerServer.
func (k *KafkaServer) Publish(topic string, body []byte, key []byte, headers map[string][]string) error {
	panic("unimplemented")
}

// Subscribe implements BrokerServer.
func (k *KafkaServer) Subscribe(topic string, handler BrokerHandler) error {
	panic("unimplemented")
}

func main() {

}
