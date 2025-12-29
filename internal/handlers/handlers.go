// Package handlers handle the http requests
package handlers

import (
	"encoding/json"
	"fmt"
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

// NewTestRepository creates new TestRepository
func NewTestRepository(app *config.AppConfig) *Repository {
	return &Repository{
		App: app,
		DB:  dbrepo.NewTestDBRepo(app),
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
		repo.App.Session.Put(r.Context(), "error", "failed to parse the start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "failed to parse the end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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
	err := r.ParseForm()
	if err != nil {
		resp := JSONResponse{
			Ok:      false,
			Message: "Internal server error",
		}
		out, _ := json.MarshalIndent(resp, "", "\n")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "failed to parse the start date")
		http.Redirect(w, r, "/search-availability", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "failed to parse the end date")
		http.Redirect(w, r, "/search-availability", http.StatusTemporaryRedirect)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "failed to parse the room id")
		http.Redirect(w, r, "/search-availability", http.StatusTemporaryRedirect)
		return
	}

	isAvailable, err := repo.DB.SearchAvailabilityByDatesByRoomID(startDate, endDate, roomID)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "failed to search availability")
		http.Redirect(w, r, "/search-availability", http.StatusTemporaryRedirect)
		return
	}

	var message string
	if isAvailable {
		message = "Available"
	} else {
		message = "Unavailable"
	}

	resp := JSONResponse{
		Ok:        isAvailable,
		Message:   message,
		RoomID:    roomID,
		StartDate: sd,
		EndDate:   ed,
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
		repo.App.Session.Put(r.Context(), "error", "cannot get the reservation from the session.\n")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := repo.DB.GetRoomByID(res.RoomID)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "cannot get the room.\n")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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
	err := r.ParseForm()
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "cannot parse the form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	form := forms.New(r.PostForm)
	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "cannot parse the start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "cannot parse the end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "invalid room id")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res := models.Reservation{}
	res.FirstName = form.Get("first_name")
	res.LastName = form.Get("last_name")
	res.Email = form.Get("email")
	res.Phone = form.Get("phone")
	res.StartDate = startDate
	res.EndDate = endDate
	res.RoomID = roomID
	res.CreatedAt = time.Now()
	res.UpdatedAt = time.Now()

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
			repo.App.Session.Put(r.Context(), "error", "cannot insert the reservation into the database")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}

		restriction := models.RoomRestriction{
			StartDate:     startDate,
			EndDate:       endDate,
			ReservationID: newReservationID,
			RoomID:        roomID,
			RestrictionID: 1,
		}

		// send email notificaions to guest
		htmlMessage := fmt.Sprintf(`
			<strong>Reservation Confirmation</strong><br>
			Dear %s:<br>
			This is confirmation from %s to %s
			`, res.FirstName, sd, ed)

		msg := models.MailData{
			To:      res.Email,
			From:    "example@example.com",
			Subject: "Reservation Confirmation",
			Content: htmlMessage,
		}

		repo.App.MailChan <- msg

		// send email to property owner
		htmlMessage = fmt.Sprintf(`
			<strong>Reservation Notification</strong><br>
			A reservation has been made for from %s to %s
			`, sd, ed)

		msg = models.MailData{
			To:      "owner@example.com",
			From:    "example@example.com",
			Subject: "Reservation Notification",
			Content: htmlMessage,
		}

		repo.App.MailChan <- msg

		_, err = repo.DB.InsertRoomRestriction(restriction)
		if err != nil {
			repo.App.Session.Put(r.Context(), "error", "cannot insert the room restriction")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
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

	layout := "2006-01-02"
	sd := reservation.StartDate.Format(layout)
	ed := reservation.EndDate.Format(layout)
	stringMap := map[string]string{}
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

func (repo *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	ID, err := strconv.Atoi(r.URL.Query().Get("room_id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	sd := r.URL.Query().Get("start")
	ed := r.URL.Query().Get("end")

	layout := "2006-01-02"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var res models.Reservation
	room, err := repo.DB.GetRoomByID(ID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.Room.RoomName = room.RoomName
	res.StartDate = startDate
	res.EndDate = endDate
	res.RoomID = ID

	repo.App.Session.Put(r.Context(), "reservation", res)
	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

type JSONResponse struct {
	Ok        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    int    `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}
