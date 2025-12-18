package render

import (
	"booking-app/internal/config"
	"booking-app/internal/models"
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
)

var testApp config.AppConfig
var session *scs.SessionManager

func TestMain(m *testing.M) {

	// change this to true when in production
	testApp.InProduction = false

	// setting up gob to handler complex data type
	gob.Register(models.Reservation{})

	// setting up session manager
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	// setting session to app config
	testApp.Session = session
	app = &testApp

	os.Exit(m.Run())
}
