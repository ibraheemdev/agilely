package validator

import (
	"log"
	"testing"
)

func TestValidate(t *testing.T) {
	v := &Validator{}
	v.Validate(true == false, "true should equal true")
	if v.Errors == nil {
		t.Error("expected validation error")
	}
}

func TestStringify(t *testing.T) {
	v := &Validator{}
	v.ValidatePresenceOf("name", "")
	v.ValidatePresenceOf("name", "")
	msgs := Stringify(v.Errors)
	if len(msgs) != 2 {
		t.Errorf("expected 2 validation errors, got %v", len(msgs))
	}
}

func TestValidatePresenceOf(t *testing.T) {
	v := &Validator{}
	v.ValidatePresenceOf("name", "")
	if v.Errors == nil {
		t.Error("expected validation error for blank field")
	}
}

func TestValidateMaxLengthOf(t *testing.T) {
	v := &Validator{}
	v.ValidateMaxLengthOf("name", "aaa", 2)
	log.Println(v.Errors)
	if v.Errors == nil {
		t.Error("expected validation error for long field")
	}
}

func TestValidateMinLengthOf(t *testing.T) {
	v := &Validator{}
	v.ValidateMinLengthOf("name", "a", 2)
	if v.Errors == nil {
		t.Error("expected validation error for short field")
	}
}
