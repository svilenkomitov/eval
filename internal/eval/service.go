package eval

import (
	"errors"
	"strconv"
	"strings"
)

type Service interface {
	Evaluate(expression string) (int, error)
}

type service struct {
}

func New() Service {
	return &service{}
}

func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func (s service) Evaluate(expression string) (int, error) {
	expression = strings.ToLower(strings.TrimSpace(expression))
	var expressionPrefix = strings.ToLower("What is")
	var expressionSuffix = strings.ToLower("?")

	if !strings.HasPrefix(expression, expressionPrefix) ||
		!strings.HasSuffix(expression, expressionSuffix) {
		return 0, errors.New("invalid question")
	}

	expression = strings.TrimSpace(expression[len(expressionPrefix) : len(expression)-len(expressionSuffix)])

	elements := strings.Fields(expression)
	if len(elements) == 0 {
		return 0, errors.New("invalid expression")
	}

	result, err := strconv.Atoi(elements[0])
	if err != nil {
		return 0, errors.New("expected number")
	}
	elements = elements[1:]

	result, err = eval(result, elements)
	if err != nil {
		return 0, err
	}

	return result, nil
}

func eval(result int, elements []string) (int, error) {
	if len(elements) == 0 {
		return result, nil
	}

	x, operation, _ := parseNext(elements)
	result, err := calculate(result, x, operation)
	if err != nil {
		return 0, err
	}

	//TODO: fix this
	if operation == DividedBy || operation == MultipliedBy {
		return eval(result, elements[3:])
	}
	return eval(result, elements[2:])
}

func calculate(x int, y int, operation Operation) (int, error) {
	switch operation {
	case Plus:
		return x + y, nil
	case Minus:
		return x - y, nil
	case MultipliedBy:
		return x * y, nil
	case DividedBy:
		return x / y, nil //TODO: divide by zero check
	default:
		return 0, errors.New("unsupported operation")
	}
}

func parseNext(elements []string) (int, Operation, error) {
	if len(elements) < 2 {
		return 0, "", errors.New("invalid expression")
	}

	//TODO: fix this
	opr := elements[0]
	numberIdx := 1
	if elements[0] == "divided" || elements[0] == "multiplied" {
		opr = strings.Join([]string{elements[0], elements[1]}, " ")
		numberIdx = 2
	}

	operation, isValid := ToOperation(opr)
	if !isValid {
		return 0, "", errors.New("invalid operation")
	}

	x, err := strconv.Atoi(elements[numberIdx])
	if err != nil {
		return 0, "", errors.New("expected number")
	}
	return x, operation, nil
}
