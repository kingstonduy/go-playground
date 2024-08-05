package utils

import "time"

func NewStringPointer(s string) *string {
	return &s
}

func NewIntPointer(i int) *int {
	return &i
}

func NewInt64Pointer(i int64) *int64 {
	return &i
}

func NewFloat64Pointer(i float64) *float64 {
	return &i
}

func NewTimePointer(t time.Time) *time.Time {
	return &t
}
