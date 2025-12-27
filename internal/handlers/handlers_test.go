package handlers

import (
	"context"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"booking-app/internal/models"
)

type postData struct {
	key   string
	value string
}

var theTests = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"gq", "/generals-quarters", "GET", http.StatusOK},
	{"ms", "/majors-suite", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"make-reservation", "/make-reservation", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	// {"search-availability", "/search-availability", "POST", []postData{
	// 	{
	// 		value: "2025-02-01",
	// 		key:   "start",
	// 	},
	// 	{
	// 		value: "2025-04-01",
	// 		key:   "end",
	// 	},
	// }, http.StatusOK},
	// {"make-reservation", "/make-reservation", "POST", []postData{
	// 	{
	// 		value: "Khadga",
	// 		key:   "first_name",
	// 	},
	// 	{
	// 		value: "Shrestha",
	// 		key:   "last_name",
	// 	},
	// 	{
	// 		value: "example@example.com",
	// 		key:   "email",
	// 	},
	// 	{
	// 		value: "9812345678",
	// 		key:   "phone",
	// 	},
	// }, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTests {
		if e.method == "GET" {
			res, err := ts.Client().Get(ts.URL + e.url)
			if err != nil {
				log.Fatal(err)
			}

			if res.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s expected %d but got %d status code,\n", e.name, e.expectedStatusCode, res.StatusCode)
			}
		}

		// else {
		// 	values := url.Values{}
		// 	for _, x := range e.params {
		// 		values.Add(x.key, x.value)
		// 	}
		//
		// 	res, err := ts.Client().PostForm(ts.URL+e.url, values)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		//
		// 	if res.StatusCode != e.expectedStatusCode {
		// 		t.Errorf("for %s expected %d but got %d status code,\n", e.name, e.expectedStatusCode, res.StatusCode)
		// 	}
		// }
	}
}

func TestRepository_Reservation(t *testing.T) {
	reservation := models.Reservation{
		ID:     1,
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req := httptest.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	// put the reservation into the session
	session.Put(ctx, "reservation", reservation)

	w := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong status code expected %d but got %d.\n", w.Code, http.StatusOK)
	}

	// request without the reservation in the session
	req = httptest.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	w = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong status code expected %d but got %d.\n", w.Code, http.StatusOK)
	}

	// request with the room  that is not exist
	req = httptest.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	// put the reservation into the session
	reservation.RoomID = 12
	session.Put(ctx, "reservation", reservation)

	w = httptest.NewRecorder()
	handler = http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong status code expected %d but got %d.\n", w.Code, http.StatusOK)
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}
	return ctx
}
