// Package handlers handle the http requests
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"booking-app/internal/config"
	"booking-app/internal/driver"
	"booking-app/internal/forms"
	"booking-app/internal/helpers"
	"booking-app/internal/models"
	"booking-app/internal/render"
	"booking-app/internal/repository"
	"booking-app/internal/repository/dbrepo"
)

var Repo *Repository

// Repository holds the Repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// NewRepository creates new NewRepository
func NewRepository(app *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: app,
		DB:  dbrepo.NewPostgresDBRepo(db.SQL, app),
	}
}

// NewHandler sets the repos
func NewHandler(r *Repository) {
	Repo = r
}

func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	startDate := r.Form.Get("start")
	endDate := r.Form.Get("end")

	fmt.Println("start date:", startDate)
	fmt.Println("end date:", endDate)

	w.Write([]byte("posted to search availability..."))
}

func (repo *Repository) PostAvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := JSONResponse{
		Ok:      true,
		Message: "Available",
	}

	out, err := json.MarshalIndent(resp, "", "\n")
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

func (repo *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]any)
	data["reservation"] = emptyReservation

	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

func (repo *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		log.Println(err)
		return
	}

	form := forms.New(r.PostForm)

	reservation := models.Reservation{
		FirstName: form.Get("first_name"),
		LastName:  form.Get("last_name"),
		Email:     form.Get("email"),
		Phone:     form.Get("phone"),
	}

	// check required fields
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)
	form.MinLength("last_name", 3, r)
	form.IsEmail("email")

	// check if the form is valid or not
	isValid := form.Valid()

	if isValid {
		repo.App.Session.Put(r.Context(), "reservation", reservation)
		http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
	} else {
		data := make(map[string]any)
		data["reservation"] = reservation

		render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
	}
}

func (repo *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := repo.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		repo.App.ErrorLog.Println("Cannot get the reservation from session")
		repo.App.Session.Put(r.Context(), "error", "cannot get the reservation-summary")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	// remove the reservation from session
	repo.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]any)
	data["reservation"] = reservation
	render.RenderTemplate(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

type JSONResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}
