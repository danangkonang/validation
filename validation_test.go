package validation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequired(t *testing.T) {
	type T struct {
		Required string `json:"required" validate:"required"`
	}
	models := []struct {
		name     string
		data     T
		expected []*validationErrors
	}{
		{
			name: "required ok",
			data: T{
				Required: "foo",
			},
		},
		{
			name: "required bad",
			data: T{
				Required: "",
			},
			expected: []*validationErrors{
				{
					Key:     "required",
					Message: "required",
				},
			},
		},
	}
	for _, m := range models {
		t.Run(m.name, func(t *testing.T) {
			r, err := MustValid(m.data)
			if err != nil {
				require.Equal(t, m.expected, r)
			} else {
				require.Equal(t, nil, err)
			}
		})
	}
}

func TestEmail(t *testing.T) {
	type T struct {
		Value string `json:"email" validate:"email"`
	}
	models := []struct {
		name     string
		data     T
		expected []*validationErrors
	}{
		{
			name: "email ok",
			data: T{
				Value: "foo@email.com",
			},
		},
		{
			name: "empty email",
			data: T{
				Value: "",
			},
			expected: []*validationErrors{
				{
					Key:     "email",
					Message: "email",
				},
			},
		},
		{
			name: "email bad 1",
			data: T{
				Value: "email",
			},
			expected: []*validationErrors{
				{
					Key:     "email",
					Message: "email",
				},
			},
		},
		{
			name: "email bad 2",
			data: T{
				Value: "email@",
			},
			expected: []*validationErrors{
				{
					Key:     "email",
					Message: "email",
				},
			},
		},
		{
			name: "email bad 3",
			data: T{
				Value: "@email",
			},
			expected: []*validationErrors{
				{
					Key:     "email",
					Message: "email",
				},
			},
		},
		{
			name: "email bad 4",
			data: T{
				Value: "@email.com",
			},
			expected: []*validationErrors{
				{
					Key:     "email",
					Message: "email",
				},
			},
		},
		{
			name: "email bad 5",
			data: T{
				Value: "email.com",
			},
			expected: []*validationErrors{
				{
					Key:     "email",
					Message: "email",
				},
			},
		},
	}
	for _, m := range models {
		t.Run(m.name, func(t *testing.T) {
			r, err := MustValid(m.data)
			if err != nil {
				require.Equal(t, m.expected, r)
			} else {
				require.Equal(t, nil, err)
			}
		})
	}
}

// func TestValidation(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		rule     string
// 		data     string
// 		axpected interface{}
// 	}{
// 		{
// 			name:     "required",
// 			rule:     "required",
// 			data:     "data",
// 			axpected: nil,
// 		},
// 		{
// 			name:     "minlength",
// 			rule:     "minlength:1",
// 			data:     "satu",
// 			axpected: nil,
// 		},
// 		{
// 			name:     "maxlength",
// 			rule:     "maxlength:5",
// 			data:     "lima",
// 			axpected: nil,
// 		},
// 		{
// 			name:     "email",
// 			rule:     "email",
// 			data:     "example@email.com",
// 			axpected: nil,
// 		},
// 	}
// 	for _, test := range tests {
// 		t.Run(test.name, func(t *testing.T) {
// 			result := Validation(test.data, test.rule)
// 			require.Equal(t, test.axpected, result)
// 		})
// 	}
// }
