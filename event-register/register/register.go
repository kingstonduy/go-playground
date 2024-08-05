package register

import (
	"encoding/json"
	"fmt"
	"go-playground/event-register/domain"
)

type EventRegister struct {
	m map[string]handlerEntry
}

type handlerEntry struct {
	handlerFunc func(interface{}) error
	eventType   interface{}
}

func NewEventRegister() EventRegister {
	// specify the size of map firsrt, it will automatically allocate later
	return EventRegister{
		m: make(map[string]handlerEntry),
	}
}

func (r *EventRegister) RegisterEvent(commandType string, eventType interface{}, handlerFunc func(interface{}) error) {
	r.m[commandType] = handlerEntry{
		handlerFunc: handlerFunc,
		eventType:   eventType,
	}
}

// HandleEvent looks up and executes the handler for an event
func (r *EventRegister) HandleEvent(eventType string, payloadStruct interface{}, event domain.Event) error {
	fmt.Println(event.Payload)
	err := json.Unmarshal([]byte(event.Payload), &payloadStruct)
	if err != nil {
		panic(err)
	}
	handlerFunc := r.m[eventType].handlerFunc

	fmt.Println(payloadStruct)
	return handlerFunc(payloadStruct)
}

func CastingEvent[T any](commandType string, payload string) T {
	var payloadStruct T
	json.Unmarshal([]byte(payload), &payload)
	return payloadStruct
}
