package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/random", nil)
	form := New(r.PostForm)
	isValid := form.Valid()

	if !isValid {
		t.Error("got invalid when it should be valid")
	}
}

func TestForm_Required(t *testing.T) {
	r := httptest.NewRequest("POST", "/random", nil)
	form := New(r.PostForm)

	form.Required("name", "email")
	isValid := form.Valid()

	if isValid {
		t.Error("got form valid when required fields are missing...")
	}

	data := url.Values{}
	data.Add("name", "khadga shrestha")
	data.Add("email", "example@example.com")
	r = httptest.NewRequest("POST", "/random", nil)
	r.PostForm = data

	form = New(r.PostForm)
	form.Required("name", "email")
	isValid = form.Valid()

	if !isValid {
		t.Error("got form invalid when it should be valid")
	}

}

func TestForm_MinLength(t *testing.T) {
	r := httptest.NewRequest("POST", "/random", nil)
	data := url.Values{}
	data.Add("name", "ku")
	r.PostForm = data

	form := New(r.PostForm)
	form.MinLength("name", 3, r)
	isValid := form.Valid()

	if isValid {
		t.Error("got form valid when minlenght is not valid")
	}

	data = url.Values{}
	data.Add("name", "khadga shrestha")
	r.PostForm = data

	form = New(r.PostForm)
	isValid = form.Valid()

	if !isValid {
		t.Error("got form invalid when it should be valid")
	}

}

func TestForm_Has(t *testing.T) {
	r := httptest.NewRequest("POST", "/random", nil)
	data := url.Values{}
	r.PostForm = data

	form := New(r.PostForm)
	form.Has("name")
	isValid := form.Valid()

	if isValid {
		t.Error("got form valid when name has no value")
	}

}
