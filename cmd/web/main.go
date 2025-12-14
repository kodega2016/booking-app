package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"booking-app/pkg/config"
	"booking-app/pkg/handlers"
	"booking-app/pkg/render"

	"github.com/alexedwards/scs/v2"
)

const port = 8080

var app config.AppConfig
var session *scs.SessionManager

func main() {

	// change this to true when in production
	app.InProduction = false

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
	}

	app.TemplateCache = tc
	app.UseCache = false

	render.NewRenderTemplate(&app)

	repo := handlers.NewRepository(&app)
	handlers.NewHandler(repo)

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
