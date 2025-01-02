package validation

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRequired(t *testing.T) {
	type T struct {
		Data string `json:"data" validate:"required"`
	}
	models := []struct {
		name     string
		data     T
		expected []*ValidationErrorMessage
	}{
		{
			name: "required ok",
			data: T{
				Data: "foo",
			},
		},
		{
			name: "required bad",
			data: T{
				Data: "",
			},
			expected: []*ValidationErrorMessage{
				{
					Field:   "data",
					Message: []interface{}{"This field is required"},
				},
			},
		},
	}
	v := New()
	for _, m := range models {
		t.Run(m.name, func(t *testing.T) {
			r, err := v.Validate(m.data)
			if err != nil {
				require.Equal(t, m.expected, r)
			} else {
				require.Equal(t, nil, err)
			}
		})
	}
}

func TestSubRequired(t *testing.T) {
	type Children struct {
		Name   string `json:"name" validate:"required"`
		Active bool   `json:"active" validate:"required"`
		Level  string `json:"level" enum:"beginner,intermediate,advanced"`
	}
	type Parent struct {
		UserId   int64      `json:"user_id" validate:"required"`
		Children []Children `json:"skill" validate:"required"`
	}
	models := []struct {
		name     string
		data     Parent
		expected []*ValidationErrorMessage
	}{
		{
			name: "required ok",
			data: Parent{
				UserId: 1,
				Children: []Children{
					{
						Name:   "foo",
						Active: true,
						Level:  "beginner",
					},
				},
			},
		},
	}
	v := New()
	for _, m := range models {
		t.Run(m.name, func(t *testing.T) {
			r, err := v.Validate(m.data)
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
		Data string `json:"email" validate:"email"`
	}
	models := []struct {
		name     string
		data     T
		expected []*ValidationErrorMessage
	}{
		{
			name: "email ok",
			data: T{
				Data: "foo@email.com",
			},
		},
		{
			name: "empty email",
			data: T{
				Data: "",
			},
			expected: []*ValidationErrorMessage{
				{
					Field:   "email",
					Message: []interface{}{"This field must be a valid email address"},
				},
			},
		},
		{
			name: "email bad 1",
			data: T{
				Data: "email",
			},
			expected: []*ValidationErrorMessage{
				{
					Field:   "email",
					Message: []interface{}{"This field must be a valid email address"},
				},
			},
		},
		{
			name: "email bad 2",
			data: T{
				Data: "email@",
			},
			expected: []*ValidationErrorMessage{
				{
					Field:   "email",
					Message: []interface{}{"This field must be a valid email address"},
				},
			},
		},
		{
			name: "email bad 3",
			data: T{
				Data: "@email",
			},
			expected: []*ValidationErrorMessage{
				{
					Field:   "email",
					Message: []interface{}{"This field must be a valid email address"},
				},
			},
		},
		{
			name: "email bad 4",
			data: T{
				Data: "@email.com",
			},
			expected: []*ValidationErrorMessage{
				{
					Field:   "email",
					Message: []interface{}{"This field must be a valid email address"},
				},
			},
		},
		{
			name: "email bad 5",
			data: T{
				Data: "email.com",
			},
			expected: []*ValidationErrorMessage{
				{
					Field:   "email",
					Message: []interface{}{"This field must be a valid email address"},
				},
			},
		},
	}
	v := New()
	for _, m := range models {
		t.Run(m.name, func(t *testing.T) {
			r, err := v.Validate(m.data)
			if err != nil {
				require.Equal(t, m.expected, r)
			} else {
				require.Equal(t, nil, err)
			}
		})
	}
}
