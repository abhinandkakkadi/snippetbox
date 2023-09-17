package validator

import (
	"strings"
	"unicode/utf8"
)

// Validation helpers

type Validator struct {
	FieldErrors map[string]string
}

func (v *Validator) Valid() bool {
	return len(v.FieldErrors) == 0
}

// Add a new Field error
func (v *Validator) AddFieldError(key, message string) {

	if v.FieldErrors == nil {
		v.FieldErrors = make(map[string]string)
	}

	if _, exists := v.FieldErrors[key]; !exists {
		v.FieldErrors[key] = message
	}
}

// add an error message to the FieldErrors only if a validation is not ok
func (v *Validator) CheckField(ok bool, key, message string) {
	if !ok {
		v.AddFieldError(key, message)
	}
}

// returns true if a value is not an empty string
func NotBlank(value string) bool {
	return strings.TrimSpace(value) != ""
}

// returns true if a value contains no more than n characters
func MaxChars(value string, n int) bool {
	return utf8.RuneCountInString(value) <= n
}

// returns true if a value is in a list of permitted integers
func PermittedInt(value int, permittedValues ...int) bool {
	for i := range permittedValues {
		if value == permittedValues[i] {
			return true
		}
	}
	return false
}
