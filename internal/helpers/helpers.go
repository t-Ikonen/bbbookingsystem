package helpers

import "github.com/t-Ikonen/bbbookingsystem/internal/config"

var app *config.AppConfig

//NewHelpers sets up cofig for helpers
func NewHelpers(a *config.AppConfig) {
	app = a
}
