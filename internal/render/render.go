// Package render will handles template rendering
package render

import (
	"bytes"
	"html/template"
	"log"
	"net/http"
	"path/filepath"

	"booking-app/internal/config"
	"booking-app/internal/models"

	"github.com/justinas/nosurf"
)

var app *config.AppConfig

func NewRenderTemplate(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	return td
}

func RenderTemplate(w http.ResponseWriter, r *http.Request, templ string, td *models.TemplateData) {
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
	td = AddDefaultData(td, r)
	err := t.Execute(buff, td)
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
