// Package config holds the application wide config
package config

import (
	"html/template"
	"log"

	"booking-app/internal/models"

	"github.com/alexedwards/scs/v2"
)

// AppConfig holds the application wide config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}
