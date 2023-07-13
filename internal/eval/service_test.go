package eval

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestService_Evaluate(t *testing.T) {
	service := New()
	t.Run("only one element, no operations", func(t *testing.T) {
		result, err := service.Evaluate("What is 5?")
		assert.NoError(t, err)
		assert.Equal(t, 5, result)
	})

	t.Run("single operation", func(t *testing.T) {
		testCase := []struct {
			expression string
			result     int
		}{
			{"What is 5 plus 13?", 18},
			{"What is 7 minus 5?", 2},
			{"What is 6 multiplied by 4?", 24},
			{"What is 25 divided by 5?", 5},
		}

		for i := range testCase {
			tt := testCase[i]
			t.Run(tt.expression, func(t *testing.T) {
				result, err := service.Evaluate(tt.expression)
				assert.NoError(t, err)
				assert.Equal(t, tt.result, result)
			})
		}
	})

	t.Run("multiple operations", func(t *testing.T) {
		testCase := []struct {
			expression string
			result     int
		}{
			{"What is 3 plus 2 multiplied by 3?", 15},
		}

		for i := range testCase {
			tt := testCase[i]
			t.Run(tt.expression, func(t *testing.T) {
				result, err := service.Evaluate(tt.expression)
				assert.NoError(t, err)
				assert.Equal(t, tt.result, result)
			})
		}
	})

	t.Run("unsupported operation", func(t *testing.T) {
		result, err := service.Evaluate("What is 52 cubed?")
		assert.ErrorIs(t, err, NewUnsupportedOperationError("cubed"))
		assert.Equal(t, 0, result)
	})

	t.Run("unsupported operation", func(t *testing.T) {
		expression := "Who is the President of the United States?"
		result, err := service.Evaluate(expression)
		assert.ErrorIs(t, err, NewInvalidQuestionError(expression))
		assert.Equal(t, 0, result)
	})

	t.Run("invalid syntax", func(t *testing.T) {
		expression := "What is 1 plus plus 2?"
		result, err := service.Evaluate(expression)
		assert.ErrorIs(t, err, NewInvalidSyntaxError(expression))
		assert.Equal(t, 0, result)
	})
}
