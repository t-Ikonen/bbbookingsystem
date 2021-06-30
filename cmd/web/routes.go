package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/t-Ikonen/bbbookingsystem/internal/config"
	"github.com/t-Ikonen/bbbookingsystem/internal/handlers"
)

func Routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	//mux.Use(WriteToConsole)
	mux.Use(NoSurf)
	mux.Use(session.LoadAndSave)
	//routes
	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/contact", handlers.Repo.Contact)
	mux.Get("/snowsuite", handlers.Repo.Snowsuite)

	mux.Get("/booking", handlers.Repo.Booking)
	mux.Post("/booking", handlers.Repo.PostBooking)
	mux.Post("/bookingjson", handlers.Repo.BookingJSON)

	mux.Get("/frostsuite", handlers.Repo.Frostsuite)
	mux.Get("/northernlights", handlers.Repo.Northernlights)
	mux.Get("/reservation", handlers.Repo.Reservation)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
