package validation

import (
	"strings"
	"sync"
	"unicode/utf8"
)

type Validator struct {
	once        sync.Once
	FieldErrors map[string]string
}

func (v *Validator) initialize() {
	v.once.Do(func() {
		v.FieldErrors = make(map[string]string)
	})
}

// Valid returns true if the FieldErrors map doesn't contain any entries.
func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// AddFieldError adds an error message to the FieldErrors map (so long as no
// entry already exists for the given key).
func (v *Validator) addFieldError(key, message string) {
	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// CheckField adds an error message to the FieldErrors map only if a
// validation check is not 'ok'.
func (v *Validator) CheckField(ok bool, key, message string) {
	v.initialize()
	if !ok {
		v.addFieldError(key, message)
	}
}

// NotBlank returns true if a value is not an empty string.
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// MaxChars returns true if a value contains no more than n characters.
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// PermittedInt returns true if a value is in a list of permitted integers.
func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}
