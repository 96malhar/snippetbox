package validation

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNotBlank(t *testing.T) {
	testcases := []struct {
		name  string
		input string
		want  bool
	}{
		{name: "Valid", input: "abcde", want: true},
		{name: "Invalid", input: "", want: false},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := NotBlank(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestMaxChars(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		maxChars int
		want     bool
	}{
		{name: "Valid", input: "abcde", maxChars: 10, want: true},
		{name: "Invalid", input: "abcde", maxChars: 3, want: false},
		{name: "Equal to max chars", input: "abcde", maxChars: 5, want: true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := MaxChars(tc.input, tc.maxChars)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestMinChars(t *testing.T) {
	testcases := []struct {
		name     string
		input    string
		minChars int
		want     bool
	}{
		{name: "Valid", input: "abcde", minChars: 3, want: true},
		{name: "Invalid", input: "abcde", minChars: 10, want: false},
		{name: "Equal to min chars", input: "abcde", minChars: 5, want: true},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := MinChars(tc.input, tc.minChars)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestMatches_ValidEmail(t *testing.T) {
	testcases := []struct {
		input string
	}{
		{input: "abc@gmail.com"},
		{input: "123@gmail.com"},
		{input: "123abc@gmail.com"},
		{input: "ab12cd@yahoo.com"},
		{input: "ab12.cd@yahoo.com"},
	}

	for _, tc := range testcases {
		t.Run(tc.input, func(t *testing.T) {
			got := Matches(tc.input, EmailRX)
			assert.True(t, got)
		})
	}
}

func TestMatches_InvalidEmail(t *testing.T) {
	testcases := []struct {
		input string
	}{
		{input: "abcgmail.com"},
		{input: "abc.gmail.com"},
		{input: "abc@gmail..com"},
		{input: "abc.com"},
		{input: "123@gmail,com"},
		{input: "123,abc@gmail.com"},
		{input: "ab12,cd@yahoo.com"},
		{input: "@yahoo.com"},
	}

	for _, tc := range testcases {
		t.Run(tc.input, func(t *testing.T) {
			got := Matches(tc.input, EmailRX)
			assert.False(t, got)
		})
	}
}

func TestPermittedValue(t *testing.T) {
	testcases := []struct {
		name            string
		input           int
		permittedValues []int
		want            bool
	}{
		{name: "Valid", input: 2, permittedValues: []int{2, 4, 6, 8}, want: true},
		{name: "Invalid", input: 2, permittedValues: []int{4, 6, 8}, want: false},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := PermittedValue(tc.input, tc.permittedValues...)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestValidator_CheckField(t *testing.T) {
	v := &Validator{}
	v.initialize()

	v.CheckField(true, "fieldKey1", "errorMessage1")
	v.CheckField(false, "fieldKey2", "errorMessage2")

	assert.NotContains(t, "fieldKey1", v.FieldErrors)
	assert.Equal(t, "errorMessage2", v.FieldErrors["fieldKey2"])
}

func TestValidator_CheckNonField(t *testing.T) {
	v := &Validator{}
	v.initialize()

	v.CheckNonField(true, "errorMessage1")
	v.CheckNonField(false, "errorMessage2")

	assert.NotContains(t, v.NonFieldErrors, "errorMessage1")
	assert.Contains(t, v.NonFieldErrors, "errorMessage2")
}

func TestValidator_Valid(t *testing.T) {
	testcases := []struct {
		name      string
		validator *Validator
		wantValid bool
	}{
		{
			name:      "Is Valid",
			validator: &Validator{},
			wantValid: true,
		},
		{
			name: "Not valid with field errors",
			validator: &Validator{
				FieldErrors: map[string]string{
					"key1": "Error1",
					"Key2": "Error2",
				},
			},
			wantValid: false,
		},
		{
			name: "Not valid with non-field errors",
			validator: &Validator{
				NonFieldErrors: []string{"Error1", "Error2"},
			},
			wantValid: false,
		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			gotValid := tc.validator.Valid()
			assert.Equal(t, tc.wantValid, gotValid)
		})
	}
}
