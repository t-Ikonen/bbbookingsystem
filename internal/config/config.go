package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
	"github.com/t-Ikonen/bbbookingsystem/internal/models"
)

//Appcongi is configuration stuct for the app
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	MailChan      chan models.MailData
}
