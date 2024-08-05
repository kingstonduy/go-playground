package usecase

import (
	"fmt"
	"go-playground/event-register/domain"
)

func HandleFooEvent(payload *domain.Foo) error {
	fmt.Println("The Message of Foo event is ", payload.Message)
	return nil
}

func HandleBarEvent(payload *domain.Bar) error {
	fmt.Println("The Value of Bar event is ", payload.Value)
	return nil
}

func HandleNopEvent(payload *domain.Nop) error {
	fmt.Println("The Result of Nop event is ", payload.Result)
	return nil
}
