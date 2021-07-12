package main

import (
	"fmt"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/t-Ikonen/bbbookingsystem/internal/config"
)

func TestRoutes(t *testing.T) {

	var app config.AppConfig

	mux := Routes(&app)
	switch v := mux.(type) {
	case *chi.Mux:
		//ok, do nothing
	default:
		t.Error(fmt.Sprintf("Type is not chi.Mux in Routes(), it is of type %T", v))
	}
}
