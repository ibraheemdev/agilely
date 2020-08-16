package validator

import (
	"fmt"
	"strings"
)

// Validator : struct field validations
type Validator struct {
	Errors []error
}

// ValidationErrors : A string of error messages
type ValidationErrors []string

// Validate :
func (v *Validator) Validate(valid bool, msg string, args ...interface{}) {
	if !valid {
		v.Errors = append(v.Errors, fmt.Errorf(msg, args...))
	}
}

// Stringify : Stringifies a slice of errors
func Stringify(errs []error) []string {
	strErrors := make([]string, len(errs))
	for i, err := range errs {
		strErrors[i] = err.Error()
	}
	return strErrors
}

// ValidatePresenceOf : validates presence of struct string field
func (v *Validator) ValidatePresenceOf(fieldName string, fieldValue string) {
	cond := len(strings.TrimSpace(fieldValue)) > 0
	v.Validate(cond, "%s cannot be blank", fieldName)
}

// ValidateMaxLengthOf : validates maximum character length of struct string field
func (v *Validator) ValidateMaxLengthOf(fieldName string, fieldValue string, max int) {
	cond := len(fieldValue) < max
	v.Validate(cond, "%s cannot be greater than %d characters", fieldName, max)
}

// ValidateMinLengthOf : validates minimum character length of struct string field
func (v *Validator) ValidateMinLengthOf(fieldName string, fieldValue string, min int) {
	cond := len(fieldValue) > min
	v.Validate(cond, "%s must be at least %d characters", fieldName, min)
}
