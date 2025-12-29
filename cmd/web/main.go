package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"booking-app/internal/config"
	"booking-app/internal/driver"
	"booking-app/internal/handlers"
	"booking-app/internal/helpers"
	"booking-app/internal/models"
	"booking-app/internal/render"
	"booking-app/internal/repository/dbrepo"

	"github.com/alexedwards/scs/v2"
)

const port = 8080

var (
	app      config.AppConfig
	session  *scs.SessionManager
	infoLog  *log.Logger
	errorLog *log.Logger
)

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}

	// close the database connection
	defer db.SQL.Close()

	// close the mail data channel
	defer close(app.MailChan)

	// listen for email
	listenForMail()

	msg := models.MailData{
		To:      "nishuka@gmail.com",
		From:    "example@example.com",
		Subject: "Just a reminder...:)",
		Content: "I am just missing you baby....",
	}

	app.MailChan <- msg

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

func run() (*driver.DB, error) {
	// change this to true when in production
	app.InProduction = false

	// setting up gob to handler complex data type
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(models.RoomRestriction{})

	// make channel for mail data
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// setting up session manager
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	// setting session to app config
	app.Session = session

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// create template cache
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("failed to create template cache")
		return nil, err
	}

	app.TemplateCache = tc
	app.UseCache = false

	app.InfoLog.Println("connecting to the database...")
	db, err := driver.ConnectSQL("host=localhost port=5432 user=kodega password=supersecret dbname=booking sslmode=disable")
	if err != nil {
		app.ErrorLog.Println(err)
		return nil, err
	}

	app.InfoLog.Println("connected to the database")

	render.NewRenderer(&app)
	repo := handlers.NewRepository(&app, db)
	handlers.NewHandler(repo)
	dbrepo.NewPostgresDBRepo(db.SQL, &app)

	helpers.NewHelpers(&app)
	return db, nil
}
