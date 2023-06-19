package validation

import (
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
			if got != tc.want {
				t.Errorf("Got = %v; Want = %v", got, tc.want)
			}
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
			if got != tc.want {
				t.Errorf("Got = %v; Want = %v", got, tc.want)
			}
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
			if got != tc.want {
				t.Errorf("Got = %v; Want = %v", got, tc.want)
			}
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
			if got != true {
				t.Errorf("Got = %v; Want = true", got)
			}
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
			if got != false {
				t.Errorf("Got = %v; Want = false", got)
			}
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
			if got != tc.want {
				t.Errorf("Got = %v; Want = %v", got, tc.want)
			}
		})
	}
}
