package validation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidation(t *testing.T) {
	tests := []struct {
		name     string
		rule     string
		data     string
		axpected interface{}
	}{
		{
			name:     "required",
			rule:     "required",
			data:     "data",
			axpected: nil,
		},
		{
			name:     "minlength",
			rule:     "minlength:1",
			data:     "satu",
			axpected: nil,
		},
		{
			name:     "maxlength",
			rule:     "maxlength:5",
			data:     "lima",
			axpected: nil,
		},
		{
			name:     "email",
			rule:     "email",
			data:     "example@email.com",
			axpected: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := Validation(test.data, test.rule)
			require.Equal(t, test.axpected, result)
		})
	}
}
