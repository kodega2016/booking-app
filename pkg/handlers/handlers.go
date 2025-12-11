// Package handlers handle the http requests
package handlers

import (
	"net/http"

	"booking-app/pkg/render"
)

func Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, "home.page.tmpl")
}

func About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplates(w, "about.page.tmpl")
}
