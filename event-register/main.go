package main

import (
	"encoding/json"
	"go-playground/event-register/domain"
	"go-playground/event-register/register"
	"go-playground/event-register/usecase"

	"github.com/google/uuid"
)

func main() {
	foo := domain.Foo{
		Message: "duydk3",
	}
	bar := domain.Bar{
		Value: 1,
	}
	nop := domain.Nop{
		Result: "result",
	}

	fooPayload, _ := json.Marshal(foo)
	barPayload, _ := json.Marshal(bar)
	nopPayload, _ := json.Marshal(nop)

	rand_uuid := uuid.New().String()

	fooEvent := domain.Event{
		AggregateID: rand_uuid,
		CommandType: "FOO_EVENT",
		Payload:     string(fooPayload),
	}
	barEvent := domain.Event{
		AggregateID: rand_uuid,
		CommandType: "BAR_EVENT",
		Payload:     string(barPayload),
	}
	nopEvent := domain.Event{
		AggregateID: rand_uuid,
		CommandType: "NOP_EVENT",
		Payload:     string(nopPayload),
	}

	reg := register.NewEventRegister()
	reg.RegisterEvent(fooEvent.CommandType, domain.Foo{}, func(i interface{}) error {
		return usecase.HandleFooEvent(&domain.Foo{})
	})
	reg.RegisterEvent(barEvent.CommandType, domain.Bar{}, func(i interface{}) error {
		return usecase.HandleBarEvent(&domain.Bar{})
	})
	reg.RegisterEvent(nopEvent.CommandType, domain.Nop{}, func(i interface{}) error {
		return usecase.HandleNopEvent(&domain.Nop{})
	})

	reg.HandleEvent(fooEvent.CommandType, &domain.Foo{}, fooEvent)

}
