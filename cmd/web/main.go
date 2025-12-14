package main

import (
	"fmt"
	"log"
	"net/http"

	"booking-app/pkg/config"
	"booking-app/pkg/handlers"
	"booking-app/pkg/render"
)

const port = 8080

func main() {
	var app config.AppConfig

	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("failed to create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

	render.NewRenderTemplate(&app)

	repo := handlers.NewRepository(&app)
	handlers.NewHandler(repo)

	fmt.Printf("server is running on port %d\n", port)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: routes(&app),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
