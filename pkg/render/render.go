// Package render render the html templates
package render

import (
	"fmt"
	"html/template"
	"net/http"
)

func RenderTemplates(w http.ResponseWriter, tmpl string) {
	t, err := template.ParseFiles("./templates/"+tmpl, "./templates/base.layout.tmpl")
	if err != nil {
		fmt.Println("failed to parse the template")
	}
	t.Execute(w, nil)
}
