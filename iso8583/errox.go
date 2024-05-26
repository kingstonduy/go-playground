package utils_class

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

const (
	DefaultSuccessStatusCode      = 200
	DefaultSuccessResponseCode    = "00"
	DefaultSuccessResponseMessage = "Successful"

	DefaultFailureStatusCode      = 500
	DefaultFailureResponseCode    = "99"
	DefaultFailureResponseMessage = "Unknown error"
)

type Error struct {
	Status  int         `json:"status"`  // http status mapping
	Code    string      `json:"code"`    // error code
	Message string      `json:"message"` // error message
	Details interface{} `json:"details"` // error details
}

func (err *Error) Error() string {
	return err.Message
}

// Parse tries to parse a JSON string into an error. If that
// fails, it will set the given string as the error detail.
func Parse(err string) *Error {
	e := new(Error)
	errr := json.Unmarshal([]byte(err), e)
	if errr != nil {
		e.Message = err
	}
	return e
}

func NewError(status int, code string, format string, a ...interface{}) *Error {
	return &Error{
		Status:  status,
		Code:    code,
		Message: strings.TrimSpace(fmt.Sprintf(format, a...)),
	}
}

func NewErrorWithDetails(status int, code string, details interface{}, format string, a ...interface{}) *Error {
	return &Error{
		Status:  status,
		Code:    code,
		Message: strings.TrimSpace(fmt.Sprintf(format, a...)),
		Details: details,
	}
}

func UnknownError(format string, a ...interface{}) *Error {
	return NewError(
		DefaultFailureStatusCode,
		DefaultFailureResponseCode,
		fmt.Sprintf("%s: %s", DefaultFailureResponseMessage, format),
		a...)
}

//============== STANDARDS ERRORS =============================

// Failed error
func Failed(format string, a ...interface{}) error {
	return NewError(200, "01", fmt.Sprintf("Failed. %s", format), a...)
}

// Failed error with details
func FailedWithDetails(details interface{}, format string, a ...interface{}) error {
	return NewErrorWithDetails(200, "01", details, fmt.Sprintf("Failed. %s", format), a...)
}

// Validation error
func ValidationError(format string, a ...interface{}) *Error {
	return NewError(200, "02", fmt.Sprintf("Validation error. %s", format), a...)
}

// Validation error with details
func ValidationErrorWithDetails(details interface{}, format string, a ...interface{}) *Error {
	return NewErrorWithDetails(200, "02", details, fmt.Sprintf("Validation error. %s", format), a...)
}

// Not found error
func NotFound(format string, a ...interface{}) *Error {
	return NewError(200, "03", fmt.Sprintf("Not found error. %s", format), a...)
}

// Not found error with details
func NotFoundWithDetails(details interface{}, format string, a ...interface{}) error {
	return NewErrorWithDetails(200, "03", details, fmt.Sprintf("Not found error. %s", format), a...)
}

// Outbound error
func OutboundError(format string, a ...interface{}) *Error {
	return NewError(200, "04", fmt.Sprintf("Outbound error. %s", format), a...)
}

// Outbound error with details
func OutboundErrorWithDetails(details interface{}, format string, a ...interface{}) error {
	return NewErrorWithDetails(200, "04", details, fmt.Sprintf("Outbound error. %s", format), a...)
}

// Timeout error
func TimeoutError(format string, a ...interface{}) *Error {
	return NewError(200, "05", fmt.Sprintf("Timeout error. %s", format), a...)
}

// Timeout error with details
func TimeoutErrorWithDetails(details interface{}, format string, a ...interface{}) error {
	return NewErrorWithDetails(200, "05", details, fmt.Sprintf("Timeout error. %s", format), a...)
}

// Internal server error
func InternalServerError(format string, a ...interface{}) *Error {
	return NewError(500, "99", fmt.Sprintf("Internal server error. %s", format), a...)
}

// Internal server error
func InternalServerErrorWithDetails(details interface{}, format string, a ...interface{}) *Error {
	return NewErrorWithDetails(500, "99", details, fmt.Sprintf("Internal server error. %s", format), a...)
}

//============================================

func BadRequest(format string, a ...interface{}) *Error {
	return NewError(200, "400", fmt.Sprintf("Bad request error: %s", format), a...)
}

func Unauthorized(format string, a ...interface{}) *Error {
	return NewError(200, "401", fmt.Sprintf("Unauthorize error: %s", format), a...)
}

func Forbidden(format string, a ...interface{}) *Error {
	return NewError(200, "403", fmt.Sprintf("Forbidden error: %s", format), a...)
}

func MethodNotAllowed(format string, a ...interface{}) *Error {
	return NewError(200, "405", fmt.Sprintf("Method not allowed error: %s", format), a...)
}

func Timeout(format string, a ...interface{}) *Error {
	return NewError(200, "408", fmt.Sprintf("Timeout error: %s", format), a...)
}

func Conflict(format string, a ...interface{}) *Error {
	return NewError(200, "409", fmt.Sprintf("Conflict error: %s", format), a...)
}

func ToManyRequest(format string, a ...interface{}) *Error {
	return NewError(429, "429", fmt.Sprintf("Too many request error: %s", format), a...)
}

func Equal(err1 error, err2 error) bool {
	verr1, ok1 := err1.(*Error)
	verr2, ok2 := err2.(*Error)

	if ok1 != ok2 {
		return false
	}

	if !ok1 {
		return err1 == err2
	}

	if verr1.Code != verr2.Code {
		return false
	}

	return true
}

// FromError try to convert go error to *Error.
func FromError(err error) *Error {
	if err == nil {
		return nil
	}
	if verr, ok := err.(*Error); ok && verr != nil {
		return verr
	}

	return Parse(err.Error())
}

// As finds the first error in err's chain that matches *Error.
func As(err error) (*Error, bool) {
	if err == nil {
		return nil, false
	}
	var merr *Error
	if errors.As(err, &merr) {
		return merr, true
	}
	return nil, false
}
