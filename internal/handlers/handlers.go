// Package handlers handle the http requests
package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"booking-app/internal/config"
	"booking-app/internal/driver"
	"booking-app/internal/forms"
	"booking-app/internal/helpers"
	"booking-app/internal/models"
	"booking-app/internal/render"
	"booking-app/internal/repository"
	"booking-app/internal/repository/dbrepo"

	"github.com/go-chi/chi/v5"
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
	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

func (repo *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	rooms, err := repo.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		repo.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]any)
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	// store the reservation in session
	repo.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})
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
	res, ok := repo.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get the reservation from the session."))
		return
	}

	room, err := repo.DB.GetRoomByID(res.RoomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res.Room.RoomName = room.RoomName

	data := make(map[string]any)
	data["reservation"] = res

	// updating the session for Reservation
	repo.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := map[string]string{}
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

func (repo *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	res, ok := repo.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("failed to get the reservation from the session"))
	}

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		log.Println(err)
		return
	}

	form := forms.New(r.PostForm)
	res.FirstName = form.Get("first_name")
	res.LastName = form.Get("last_name")
	res.Email = form.Get("email")
	res.Phone = form.Get("phone")

	// check required fields
	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3, r)
	form.MinLength("last_name", 3, r)
	form.IsEmail("email")

	// check if the form is valid or not
	isValid := form.Valid()

	if isValid {
		newReservationID, err := repo.DB.InsertReservation(res)
		if err != nil {
			helpers.ServerError(w, err)
		}

		restriction := models.RoomRestriction{
			StartDate:     res.StartDate,
			EndDate:       res.EndDate,
			ReservationID: newReservationID,
			RoomID:        res.RoomID,
			RestrictionID: 1,
		}

		_, err = repo.DB.InsertRoomRestriction(restriction)
		if err != nil {
			helpers.ServerError(w, err)
		}

		repo.App.Session.Put(r.Context(), "reservation", res)
		http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
	} else {
		data := make(map[string]any)
		data["reservation"] = res

		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
	}
}

func (repo *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
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

	stringMap := map[string]string{}
	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

func (repo *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res, ok := repo.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, err)
		return
	}

	res.RoomID = roomID
	data := make(map[string]any)
	data["reservation"] = res

	// again put the reservation into the session
	repo.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

type JSONResponse struct {
	Ok      bool   `json:"ok"`
	Message string `json:"message"`
}
