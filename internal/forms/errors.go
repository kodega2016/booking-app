// Package forms handles the application form
package forms

// errors will hold the form errors(field:message)
type errors map[string][]string

// Add will add error message to a field
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
