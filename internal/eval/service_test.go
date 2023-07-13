package eval

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_Evaluate(t *testing.T) {
	service := New()

	testCase := []struct {
		expression string
		result     int
		err        *ResponseError
	}{
		{"What is 5?", 5, nil},
		{"What is 5 plus 13?", 18, nil},
		{"What is 7 minus 5?", 2, nil},
		{"What is 6 multiplied by 4?", 24, nil},
		{"What is 25 divided by 5?", 5, nil},
		{"What is 3 plus 2 multiplied by 3?", 15, nil},
		{"What is 52 cubed?", 0, NewUnsupportedOperationError("cubed")},
		{"Who is the President of the United States?", 0, NewInvalidQuestionError("Who is the President of the United States?")},
		{"What is 1 plus plus 2?", 0, NewInvalidSyntaxError("What is 1 plus plus 2?")},
		{"What is 5 divided by 0?", 0, NewInvalidArithmeticsError("What is 5 divided by 0?", "divide by zero")},
	}
	for i := range testCase {
		tt := testCase[i]
		t.Run(tt.expression, func(t *testing.T) {
			result, err := service.Evaluate(tt.expression)
			if tt.err == nil {
				assert.Nil(t, err)
			} else {
				assert.EqualError(t, err, tt.err.Error())
			}
			assert.Equal(t, tt.result, result)
		})
	}
}
