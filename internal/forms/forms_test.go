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

	has := form.Has("a")
	if has {
		t.Error("form shows having field when it does not")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	form = New(postedData)

	has = form.Has("a")
	if !has {
		t.Error("form shows not having field when it should")
	}
}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	form.MinLength("whatever", 5)
	if form.Valid() {
		t.Error("form shows min length for non-existent field")
	}

	isError := form.Errors.Get("whatever")
	if isError == "" {
		t.Error("should have an error, but did not get one")
	}

	postedData := url.Values{}
	postedData.Add("invalid_length", "a")
	form = New(postedData)

	form.MinLength("invalid_length", 5)
	if form.Valid() {
		t.Error("form shows min length of 5 met when data is shorter")
	}

	postedData = url.Values{}
	postedData.Add("valid_length", "long value")
	form = New(postedData)

	form.MinLength("valid_length", 5)
	if !form.Valid() {
		t.Error("form shows min length of 5 met when it is")
	}

	isError = form.Errors.Get("valid_length")
	if isError != "" {
		t.Error("should not have an error, but got one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	form := New(url.Values{})

	form.IsEmail("invalid_email")
	if form.Valid() {
		t.Error("form shows valid email for non-existent field")
	}

	postedData := url.Values{}
	postedData.Add("invalid_email", "a")
	form = New(postedData)

	form.IsEmail("invalid_email")
	if form.Valid() {
		t.Error("form shows valid email when field has invalid email")
	}

	postedData = url.Values{}
	postedData.Add("valid_email", "yyj@test.com")
	form = New(postedData)

	form.IsEmail("valid_email")
	if !form.Valid() {
		t.Error("form shows invalid email when field has valid email")
	}
}
