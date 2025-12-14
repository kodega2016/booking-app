package main

import (
	"booking-app/pkg/config"
	"booking-app/pkg/handlers"
	"net/http"

	"github.com/gorilla/pat"
)

func routes(app *config.AppConfig) http.Handler {
	mux := pat.New()
	mux.Get("/about", http.HandlerFunc(handlers.Repo.About))
	mux.Get("/", http.HandlerFunc(handlers.Repo.Home))
	return mux
}
