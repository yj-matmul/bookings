package forms

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")

	r, _ = http.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("form shows does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	hasResult := form.Has("a")
	if hasResult {
		t.Error("form shows having field when field is missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")

	r = httptest.NewRequest("POST", "/whatever", nil)
	r.PostForm = postedData
	form = New(r.PostForm)
	hasResult = form.Has("a")
	if !hasResult {
		t.Error("form shows not having field when field exist")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	postedData := url.Values{}
	postedData.Add("invalid_length", "a")
	postedData.Add("valid_length", "yyjabc")
	r.PostForm = postedData

	minLengthResult := form.MinLength("invalid_length", 3, r)
	if minLengthResult {
		t.Error("form shows valid length on field when field has invalid length")
	}

	minLengthResult = form.MinLength("valid_length", 2, r)
	if minLengthResult {
		t.Error("form shows invalid length on field when field has valid length")
	}
}

func TestForm_IsEmail(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)

	postedData := url.Values{}
	postedData.Add("invalid_email", "a")
	r.PostForm = postedData
	form := New(r.PostForm)

	form.IsEmail("invalid_email")
	if form.Valid() {
		t.Error("form shows valid email when field has invalid email")
	}

	postedData = url.Values{}
	postedData.Add("valid_email", "yyj@test.com")
	r.PostForm = postedData
	form = New(r.PostForm)

	form.IsEmail("valid_email")
	if !form.Valid() {
		t.Error("form shows invalid email when field has valid email")
	}
}
