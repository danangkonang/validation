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
		expected []*ValidationErrors
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
			expected: []*ValidationErrors{
				{
					Key: "required",
				},
			},
		},
	}
	v := New()
	for _, m := range models {
		t.Run(m.name, func(t *testing.T) {
			r, err := v.MustValid(m.data)
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
		expected []*ValidationErrors
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
			expected: []*ValidationErrors{
				{
					Key: "email",
				},
			},
		},
		{
			name: "email bad 1",
			data: T{
				Value: "email",
			},
			expected: []*ValidationErrors{
				{
					Key: "email",
				},
			},
		},
		{
			name: "email bad 2",
			data: T{
				Value: "email@",
			},
			expected: []*ValidationErrors{
				{
					Key: "email",
				},
			},
		},
		{
			name: "email bad 3",
			data: T{
				Value: "@email",
			},
			expected: []*ValidationErrors{
				{
					Key: "email",
				},
			},
		},
		{
			name: "email bad 4",
			data: T{
				Value: "@email.com",
			},
			expected: []*ValidationErrors{
				{
					Key: "email",
				},
			},
		},
		{
			name: "email bad 5",
			data: T{
				Value: "email.com",
			},
			expected: []*ValidationErrors{
				{
					Key: "email",
				},
			},
		},
	}
	v := New()
	for _, m := range models {
		t.Run(m.name, func(t *testing.T) {
			r, err := v.MustValid(m.data)
			if err != nil {
				require.Equal(t, m.expected, r)
			} else {
				require.Equal(t, nil, err)
			}
		})
	}
}
