package eval

import "fmt"

type ErrorKind string

const (
	ValidationError ErrorKind = "Validation Error"
)

type ResponseError struct {
	Kind   ErrorKind `json:"kind"`
	Reason string    `json:"reason"`
	Code   int       `json:"code"`
}

func NewResponseError(kind ErrorKind, reason string, code int) *ResponseError {
	return &ResponseError{
		Kind:   kind,
		Reason: reason,
		Code:   code,
	}
}

func (e *ResponseError) Error() string {
	return fmt.Sprintf("error [%s] occured. reason: [%s] code:%d", e.Kind, e.Reason, e.Code)
}
