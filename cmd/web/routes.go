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

	mux.Get("/chooseroom/{id}", handlers.Repo.ChooseRoom)

	mux.Get("/bookroom", handlers.Repo.BookRoom)

	mux.Get("/frostsuite", handlers.Repo.Frostsuite)
	mux.Get("/northernlights", handlers.Repo.Northernlights)

	mux.Get("/reservation", handlers.Repo.Reservation)
	mux.Post("/reservation", handlers.Repo.PostReservation)
	mux.Get("/reservationsummary", handlers.Repo.Reservationsummary)

	mux.Get("/user/login", handlers.Repo.ShowLogin)
	mux.Post("/user/login", handlers.Repo.PostShowLogin)
	mux.Get("/user/logout", handlers.Repo.ShowLogout)

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/dashboard", handlers.Repo.AdminDashboard)
		mux.Get("/statistics", handlers.Repo.AdminStatistics)

		mux.Get("/newreservations", handlers.Repo.AdminNewReservations)
		mux.Get("/allreservations", handlers.Repo.AdminAllReservations)
		mux.Get("/reservationcalendar", handlers.Repo.AdminCalendar)
	})

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
