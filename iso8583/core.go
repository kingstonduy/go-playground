package utils_class

import (
	"time"
)

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

func (r *Response[T]) SetTrace(trace Trace) {
	r.Trace = trace
}

func (r *Response[T]) GetTrace() Trace {
	return r.Trace
}

var (
	DefaultSuccessResponse = Response[interface{}]{
		Result: Result{
			StatusCode: DefaultSuccessStatusCode,
			Code:       DefaultSuccessResponseCode,
			Message:    DefaultSuccessResponseMessage,
		},
		Trace: Trace{
			Sts: time.Now().UnixMilli(),
		},
		Data: nil,
	}

	DefaultFailureResponse = Response[interface{}]{
		Result: Result{
			StatusCode: DefaultFailureStatusCode,
			Code:       DefaultFailureResponseCode,
			Message:    DefaultFailureResponseMessage,
		},
		Trace: Trace{
			Sts: time.Now().UnixMilli(),
		},
		Data: nil,
	}
)
