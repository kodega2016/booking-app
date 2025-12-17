package main

import (
	"booking-app/internal/config"
	"testing"

	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {
	var appConfig config.AppConfig
	mux := routes(&appConfig)

	switch v := mux.(type) {
	case *chi.Mux:
	default:
		t.Errorf("type is not chi.mux but we got %T", v)
	}
}
