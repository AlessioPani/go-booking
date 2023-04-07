package forms

import (
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	postData := url.Values{}
	f := New(postData)

	isValid := f.Valid()
	if !isValid {
		t.Error("form invalid when it should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	// Test with a blank form
	postData := url.Values{}
	f := New(postData)

	f.Required("required_field")
	if f.Valid() {
		t.Error("form without a required field defined as valid")
	}

	// Test with a filled form
	postData.Add("required_field1", "required_value1")
	postData.Add("required_field2", "required_value2")
	postData.Add("required_field3", "required_value3")
	f = New(postData)

	f.Required("required_field1", "required_field2", "required_field3")
	if !f.Valid() {
		t.Error("valid form with required fields defined as invalid")
	}
}

func TestForm_MinLength(t *testing.T) {
	// Test with a blank form
	postData := url.Values{}
	f := New(postData)

	f.MinLength("random_field", 10)
	if f.Valid() {
		t.Error("form validated for min lenght field in a blank form")
	}

	isError := f.Errors.Get("random_field")
	if isError == "" {
		t.Error("should have an error, but did not get one")
	}

	// Test with a filled form
	postData.Add("random_field", "test_min_length")
	f = New(postData)

	// Field ok
	f.MinLength("random_field", 10)
	if !f.Valid() {
		t.Error("form not validated with a proper valued field")
	}

	isError = f.Errors.Get("random_field")
	if isError != "" {
		t.Error("should not have an error, but did get one")
	}

	// Field not ok
	f.MinLength("random_field", 100)
	if f.Valid() {
		t.Error("form not validated with a field with less than 100 char")
	}
}

func TestForm_Has(t *testing.T) {
	// Test with a blank form
	postData := url.Values{}
	f := New(postData)

	has := f.Has("random_field")
	if has {
		t.Error("field found in a blank form")
	}

	// Test with a filled form
	postData.Add("random_field", "random_value")
	f = New(postData)

	has = f.Has("random_field")
	if !has {
		t.Error("field not found in a filled form")
	}
}

func TestForm_IsEmail(t *testing.T) {
	// Test with a blank form
	postData := url.Values{}
	f := New(postData)

	f.IsEmail("email_field")

	if f.Valid() {
		t.Error("form shows valid email in a blank form")
	}

	// Test with a filled form
	// Email field ok
	postData.Add("correct_email", "test@test.com")
	f = New(postData)

	f.IsEmail("correct_email")
	if !f.Valid() {
		t.Error("form defined invalid with a proper email value")
	}

	// Email field not ok
	postData.Add("uncorrect_email", "string")
	f = New(postData)

	f.IsEmail("uncorrect_email")
	if f.Valid() {
		t.Error("form defined valid with a wrong email value")
	}
}
