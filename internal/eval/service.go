package eval

import (
	log "github.com/sirupsen/logrus"
	"net/http"
	"strconv"
	"strings"
)

var (
	expressionPrefix = strings.ToLower("What is")
	expressionSuffix = strings.ToLower("?")
)

type Service interface {
	Evaluate(expression string) (int, *ResponseError)
	Validate(expression string) *ResponseError
	Errors() ([]Error, *ResponseError)
}

type service struct {
	repository Repository
}

func New(repository Repository) Service {
	return &service{
		repository: repository,
	}
}

func (s service) Errors() ([]Error, *ResponseError) {
	errs, err := s.repository.Errors()
	if err != nil {
		log.Error(err)
		return []Error{}, NewResponseError(
			DatabaseError,
			"failed to fetch errors",
			http.StatusInternalServerError)
	}
	return errs, nil
}

func (s service) Validate(expression string) *ResponseError {
	var respError *ResponseError
	defer func() {
		if respError != nil {
			if err := s.repository.Upsert(Error{
				Expression: expression,
				Endpoint:   Validate,
				Kind:       respError.Kind,
			}); err != nil {
				log.Error(err)
			}
		}
	}()

	isValidQuestion := func(expression string) bool {
		expression = strings.ToLower(strings.TrimSpace(expression))
		return strings.HasPrefix(expression, expressionPrefix) &&
			strings.HasSuffix(expression, expressionSuffix)
	}

	if !isValidQuestion(expression) {
		respError = NewInvalidQuestionError(expression)
		return respError
	}

	mathExpression := strings.TrimSpace(expression[len(expressionPrefix) : len(expression)-len(expressionSuffix)])
	elements := strings.Fields(mathExpression)
	if len(elements) == 0 {
		respError = NewInvalidSyntaxError(expression)
		return respError
	}

	elements, _, err := nextNumber(expression, elements)
	if err != nil {
		respError = err
		return respError
	}

	respError = validateMathExpression(expression, elements)
	return respError
}

func validateMathExpression(expression string, elements []string) *ResponseError {
	if len(elements) == 0 {
		return nil
	}

	elements, _, err := nextOperation(expression, elements)
	if err != nil {
		return err
	}

	elements, _, err = nextNumber(expression, elements)
	if err != nil {
		return err
	}

	return validateMathExpression(expression, elements)
}

func (s service) Evaluate(expression string) (int, *ResponseError) {
	var respError *ResponseError
	defer func() {
		if respError != nil {
			if err := s.repository.Upsert(Error{
				Expression: expression,
				Endpoint:   Evaluate,
				Kind:       respError.Kind,
			}); err != nil {
				log.Error(err)
			}
		}
	}()

	isValidQuestion := func(expression string) bool {
		expression = strings.ToLower(strings.TrimSpace(expression))
		return strings.HasPrefix(expression, expressionPrefix) &&
			strings.HasSuffix(expression, expressionSuffix)
	}

	if !isValidQuestion(expression) {
		respError = NewInvalidQuestionError(expression)
		return 0, respError
	}

	mathExpression := strings.TrimSpace(expression[len(expressionPrefix) : len(expression)-len(expressionSuffix)])
	elements := strings.Fields(mathExpression)
	if len(elements) == 0 {
		respError = NewInvalidSyntaxError(expression)
		return 0, respError
	}

	elements, number, err := nextNumber(expression, elements)
	if err != nil {
		respError = err
		return 0, respError
	}

	result, respError := evalMathExpression(expression, elements, number)
	return result, respError
}

func evalMathExpression(expression string, elements []string, result int) (int, *ResponseError) {
	if len(elements) == 0 {
		return result, nil
	}

	elements, operation, err := nextOperation(expression, elements)
	if err != nil {
		return 0, err
	}

	elements, number, err := nextNumber(expression, elements)
	if err != nil {
		return 0, err
	}

	result, err = calculate(expression, result, number, operation)
	if err != nil {
		return 0, err
	}

	return evalMathExpression(expression, elements, result)
}

func calculate(expression string, x int, y int, operation Operation) (int, *ResponseError) {
	switch operation {
	case Plus:
		return x + y, nil
	case Minus:
		return x - y, nil
	case MultipliedBy:
		return x * y, nil
	case DividedBy:
		if y == 0 {
			return 0, NewInvalidArithmeticsError(expression, "divide by zero")
		}
		return x / y, nil
	default:
		return 0, NewUnsupportedOperationError(string(operation))
	}
}

func nextOperation(expression string, elements []string) ([]string, Operation, *ResponseError) {
	if len(elements) == 0 {
		return elements, "", NewInvalidSyntaxError(expression)
	}

	elementsCount := 1
	operation := elements[0]
	if operation == "divided" || operation == "multiplied" {
		if len(elements) < 2 {
			return elements, "", NewInvalidSyntaxError(expression)
		}
		operation = strings.Join([]string{elements[0], elements[1]}, " ")
		elementsCount = 2
	}

	result, isValid := ToOperation(operation)
	if !isValid {
		return elements, "", NewUnsupportedOperationError(operation)
	}
	return elements[elementsCount:], result, nil
}

func nextNumber(expression string, elements []string) ([]string, int, *ResponseError) {
	if len(elements) == 0 {
		return elements, 0, NewInvalidSyntaxError(expression)
	}

	number, err := strconv.Atoi(elements[0])
	if err != nil {
		return elements, 0, NewInvalidSyntaxError(expression)
	}

	return elements[1:], number, nil
}
