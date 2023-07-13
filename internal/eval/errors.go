package eval

import (
	"fmt"
	"net/http"
)

type ErrorKind string

const (
	ValidationError      ErrorKind = "Validation Error"
	UnsupportedOperation ErrorKind = "Unsupported Operation"
	InvalidQuestion      ErrorKind = "Invalid Question"
)

type ResponseError struct {
	Kind   ErrorKind `json:"kind"`
	Reason string    `json:"reason"`
	Code   int       `json:"code"`
}

func NewResponseError(kind ErrorKind, reason string, code int) ResponseError {
	return ResponseError{
		Kind:   kind,
		Reason: reason,
		Code:   code,
	}
}

func (e ResponseError) Error() string {
	return fmt.Sprintf("error \"%s\" occured. reason: \"%s\" code:%d", e.Kind, e.Reason, e.Code)
}

type UnsupportedOperationError struct {
	ResponseError
}

func NewUnsupportedOperationError(operation string) ResponseError {
	return NewResponseError(UnsupportedOperation,
		fmt.Sprintf("operation \"%s\" is not supported. supported operations: %v",
			operation, SupportedOperations), http.StatusBadRequest,
	)
}

type InvalidQuestionError struct {
	ResponseError
}

func NewInvalidQuestionError(expression string) ResponseError {
	return NewResponseError(InvalidQuestion,
		fmt.Sprintf("question \"%s\" is not valid. valid pattern: %v",
			expression, "What is ....?"), http.StatusBadRequest,
	)
}

type InvalidSyntaxError struct {
	ResponseError
}

func NewInvalidSyntaxError(expression string) ResponseError {
	return NewResponseError(InvalidQuestion,
		fmt.Sprintf("expression \"%s\" syntax is not valid. valid syntax: %v",
			expression, "What is [number] [operation] [number] ...?"), http.StatusBadRequest,
	)
}
