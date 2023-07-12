package eval

import "strings"

type Operation string

const (
	Plus         Operation = "plus"
	Minus        Operation = "minus"
	MultipliedBy Operation = "multiplied by"
	DividedBy    Operation = "divided by"
)

var SupportedOperations = []Operation{Plus, Minus, DividedBy, MultipliedBy}

func (o Operation) ToString() string {
	return strings.ToLower(string(o))
}

func IsSupportedOperation(operation string) bool {
	for _, o := range SupportedOperations {
		if o.ToString() == operation {
			return true
		}
	}
	return false
}

func ToOperation(operation string) (Operation, bool) {
	switch strings.ToLower(operation) {
	case string(Plus):
		return Plus, true
	case string(Minus):
		return Minus, true
	case string(DividedBy):
		return DividedBy, true
	case string(MultipliedBy):
		return MultipliedBy, true
	default:
		return "", false
	}
}
