// Package render will handles template rendering
package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"booking-app/pkg/config"
)

var app *config.AppConfig

func NewRenderTemplate(a *config.AppConfig) {
	app = a
}

func RenderTemplate(w http.ResponseWriter, templ string) {
	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, _ = CreateTemplateCache()
	}

	// get template from the cache
	t, ok := tc[templ]
	if !ok {
		log.Fatal("failed to get the template cache")
	}

	// render template with tempate data
	buff := new(bytes.Buffer)
	err := t.Execute(buff, nil)
	if err != nil {
		log.Fatal(err)
	}

	_, err = buff.WriteTo(w)
	if err != nil {
		log.Fatal(err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// find the page templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// find the layout templates
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return nil, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return nil, err
			}
		}
		myCache[name] = ts
	}

	return myCache, nil
}
