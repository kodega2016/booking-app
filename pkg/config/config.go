// Package config holds the application wide config
package config

import (
	"html/template"
	"log"
)

// AppConfig holds the application wide config
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
}
