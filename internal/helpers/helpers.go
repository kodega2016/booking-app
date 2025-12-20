// Package helpers contains helper functions
package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"booking-app/internal/config"
)

var app *config.AppConfig

// NewHelpers create new helper instance
func NewHelpers(a *config.AppConfig) {
	app = a
}

func ClientError(w http.ResponseWriter, statusCode int) {
	app.InfoLog.Println("Client error with status of ", statusCode)
	http.Error(w, http.StatusText(statusCode), statusCode)
}

func ServerError(w http.ResponseWriter, error error) {
	trace := fmt.Sprintf("%s\n%s", error.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
