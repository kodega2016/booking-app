package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"booking-app/internal/config"
	"booking-app/internal/models"
	"booking-app/internal/render"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
)

var (
	app             config.AppConfig
	session         *scs.SessionManager
	pathToTemplates = "./../../templates"
	functions       = template.FuncMap{}
)

func getRoutes() http.Handler {
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
	tc, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("failed to create template cache")
	}

	app.TemplateCache = tc
	app.UseCache = true
	render.NewRenderTemplate(&app)

	mux := chi.NewRouter()
	mux.Use(session.LoadAndSave)
	// use middleware
	// mux.Use(WriteToConsole)
	// mux.Use(NoSurf)
	// mux.Use(SessionLoad)
	//
	// static file server
	fs := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fs))

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/generals-quarters", Repo.Generals)
	mux.Get("/majors-suite", Repo.Majors)
	mux.Get("/search-availability", Repo.Availability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.PostAvailabilityJSON)
	mux.Get("/make-reservation", Repo.Reservation)
	mux.Post("/make-reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)
	mux.Get("/contact", Repo.Contact)

	return mux
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// find the page templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// find the layout templates
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return nil, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return nil, err
			}
		}
		myCache[name] = ts
	}

	return myCache, nil
}
