package handlers

import (
	"encoding/gob"
	"fmt"
	"html/template"
	"log"
	"os"
	"testing"

	"net/http"
	"path/filepath"
	"time"

	"github.com/alexedwards/scs/v2"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/justinas/nosurf"
	"github.com/t-Ikonen/bbbookingsystem/internal/config"
	"github.com/t-Ikonen/bbbookingsystem/internal/models"
	"github.com/t-Ikonen/bbbookingsystem/internal/render"
)

var appCnf config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"

var functions = template.FuncMap{}

func TestMain(m *testing.M) {
	// Reservation model stored in session
	gob.Register(models.Reservation{})

	//change to true when in production
	appCnf.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	appCnf.InfoLog = infoLog

	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	appCnf.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appCnf.InProduction

	appCnf.Session = session

	mailChan := make(chan models.MailData)
	appCnf.MailChan = mailChan
	defer close(mailChan)
	listenForMail()

	tmplCache, err := CreateTestTemplateCache()
	if err != nil {
		fmt.Printf("Error crating template configuration, error %s \n", err)

		//fmt.Println(fmt.Sprintf("Error crating template configuration, error %s \n", err))
	}
	appCnf.TemplateCache = tmplCache
	appCnf.UseCache = true

	repo := NewTestRepo(&appCnf)
	NewHandlers(repo)
	render.NewRenderer(&appCnf)
	os.Exit(m.Run())
}

func listenForMail() {
	go func() {
		for {
			_ = <-appCnf.MailChan
		}
	}()
}

func getRoutes() http.Handler {

	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	//mux.Use(WriteToConsole)
	//mux.Use(NoSurf)
	mux.Use(session.LoadAndSave)
	//routes
	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/contact", Repo.Contact)
	mux.Get("/snowsuite", Repo.Snowsuite)

	mux.Get("/booking", Repo.Booking)
	mux.Post("/booking", Repo.PostBooking)
	mux.Post("/bookingjson", Repo.BookingJSON)

	mux.Get("/frostsuite", Repo.Frostsuite)
	mux.Get("/northernlights", Repo.Northernlights)

	mux.Get("/reservation", Repo.Reservation)
	mux.Post("/reservation", Repo.PostReservation)
	mux.Get("/reservationsummary", Repo.Reservationsummary)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
func NoSurf(next http.Handler) http.Handler {
	csrtHandler := nosurf.New(next)

	csrtHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   appCnf.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrtHandler
}

//SessionLoadloads and save the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl.html", pathToTemplates))

	if err != nil {

		return myCache, err
	}

	for _, page := range pages {

		name := filepath.Base(page)
		//fmt.Println("page filelistassa on", page, "ja name on ", name)

		tmplSet, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {

			return myCache, err
		}
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl.html", pathToTemplates))
		if err != nil {

			return myCache, err
		}
		if len(matches) > 0 {
			tmplSet, err = tmplSet.ParseGlob(fmt.Sprintf("%s/*.page.tmpl.html", pathToTemplates))
			if err != nil {

				return myCache, err
			}
		}
		myCache[name] = tmplSet
		//fmt.Println("tmplSet ", tmplSet)
	}

	return myCache, nil
}
