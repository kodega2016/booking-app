// Package handlers handle the http requests
package handlers

import (
	"net/http"

	"booking-app/pkg/config"
	"booking-app/pkg/render"
)

var Repo *Repository

// Repository holds the Repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepository creates new NewRepository
func NewRepository(app *config.AppConfig) *Repository {
	return &Repository{
		App: app,
	}
}

// NewHandler sets the repos
func NewHandler(r *Repository) {
	Repo = r
}

func Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "home.page.tmpl")
}

func About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "about.page.tmpl")
}
