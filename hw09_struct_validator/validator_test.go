package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type (
	UserRole      string
	NumberWrapper int
)

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string        `validate:"len:5"`
		Numbers []int         `validate:"min:10|max:50"`
		Wrap    NumberWrapper `validate:"max:100|min:25"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr ValidationErrors
	}{
		{
			in: User{
				ID:     "1",
				Name:   "Xipe-Totec",
				Age:    5,
				Email:  "hello@world",
				Role:   "superuser",
				Phones: []string{"111222333"},
				meta:   json.RawMessage{},
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "ID",
					Err:   ErrLen,
				},
				ValidationError{
					Field: "Age",
					Err:   ErrMin,
				},
				ValidationError{
					Field: "Email",
					Err:   ErrRegexp,
				},
				ValidationError{
					Field: "Role",
					Err:   ErrIn,
				},
				ValidationError{
					Field: "Phones",
					Err:   ErrLen,
				},
			},
		},
		{
			in: App{
				Version: "ololo",
				Numbers: []int{10, 15, 218},
				Wrap:    12,
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Numbers",
					Err:   ErrMax,
				},
				ValidationError{
					Field: "Wrap",
					Err:   ErrMin,
				},
			},
		},
		{
			in: Token{
				Header:    []byte("text/json"),
				Payload:   nil,
				Signature: nil,
			},
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 0,
				Body: "",
			},
			expectedErr: ValidationErrors{
				ValidationError{
					Field: "Code",
					Err:   ErrIn,
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			err := Validate(tt.in)
			assert.Equal(t, &tt.expectedErr, err)
		})
	}
}
