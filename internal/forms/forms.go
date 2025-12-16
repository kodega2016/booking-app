// Package forms handles the application form
package forms

import (
	"net/url"
	"strings"
)

type Form struct {
	url.Values
	Errors errors
}

// New initializes the form
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Has check if the form has the value in that field or not
func (f *Form) Has(field string) bool {
	x := f.Get(field)
	if x == "" {
		f.Errors.Add(field, "This field cannot be empty.")
		return false
	}
	return true
}

func (f *Form) Valid() bool {
	return len(f.Errors) == 0
}

func (f *Form) Required(fields ...string) {
	for _, field := range fields {
		value := f.Get(field)
		if strings.TrimSpace(value) == "" {
			f.Errors.Add(field, "This field cannot be blank")
		}
	}
}
