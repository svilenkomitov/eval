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
	InvalidSyntax        ErrorKind = "Invalid Syntax"
	InvalidArithmetics   ErrorKind = "Invalid Arithmetics"
	DatabaseError        ErrorKind = "Database Error"
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

func (e ResponseError) Error() string {
	return e.Reason
}

type UnsupportedOperationError struct {
	ResponseError
}

func NewUnsupportedOperationError(operation string) *ResponseError {
	return NewResponseError(UnsupportedOperation,
		fmt.Sprintf("operation [%s] is not supported. supported operations: %v",
			operation, SupportedOperations), http.StatusBadRequest,
	)
}

type InvalidQuestionError struct {
	ResponseError
}

func NewInvalidQuestionError(expression string) *ResponseError {
	return NewResponseError(InvalidQuestion,
		fmt.Sprintf("question [%s] is not valid. valid question: [%v]",
			expression, "What is ....?"), http.StatusBadRequest,
	)
}

type InvalidSyntaxError struct {
	ResponseError
}

func NewInvalidSyntaxError(expression string) *ResponseError {
	return NewResponseError(InvalidSyntax,
		fmt.Sprintf("expression [%s] syntax is not valid. valid syntax: %v",
			expression, "What is [number] [operation] [number] ...?"), http.StatusBadRequest,
	)
}

type InvalidArithmeticsError struct {
	ResponseError
}

func NewInvalidArithmeticsError(expression string, reason string) *ResponseError {
	return NewResponseError(InvalidArithmetics,
		fmt.Sprintf("expression [%s] arithmetic is not valid. reason: %v",
			expression, reason), http.StatusBadRequest,
	)
}
