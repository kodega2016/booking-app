package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"booking-app/internal/config"
	"booking-app/internal/handlers"
	"booking-app/internal/models"
	"booking-app/internal/render"

	"github.com/alexedwards/scs/v2"
)

const port = 8080

var (
	app     config.AppConfig
	session *scs.SessionManager
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
	// starting the server
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

func run() error {
	// change this to true when in production
	app.InProduction = false

	// setting up gob to handler complex data type
	gob.Register(models.Reservation{})

	// setting up session manager
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	// setting session to app config
	app.Session = session

	// create template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("failed to create template cache")
		return err
	}

	app.TemplateCache = tc
	app.UseCache = false

	render.NewRenderTemplate(&app)

	repo := handlers.NewRepository(&app)
	handlers.NewHandler(repo)
	return nil
}
