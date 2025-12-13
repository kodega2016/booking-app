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
	render.NewRenderTemplate(&app)

	http.HandleFunc("/", handlers.Home)
	http.HandleFunc("/about", handlers.About)

	fmt.Printf("server is running on port %d\n", port)
	_ = http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
