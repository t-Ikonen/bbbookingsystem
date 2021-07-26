package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/t-Ikonen/bbbookingsystem/internal/config"
	"github.com/t-Ikonen/bbbookingsystem/internal/handlers"
	"github.com/t-Ikonen/bbbookingsystem/internal/helpers"
	"github.com/t-Ikonen/bbbookingsystem/internal/models"
	"github.com/t-Ikonen/bbbookingsystem/internal/render"
)

const portNum = ":8080"

var appCnf config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

//Main of HelloWeb app
func main() {

	err := run()
	if err != nil {
		log.Fatal(err)
	}
	srv := &http.Server{
		Addr:    portNum,
		Handler: Routes(&appCnf),
	}
	fmt.Printf("Starting app on port %s for your pleasure \n", portNum)
	//fmt.Println(fmt.Sprintf("Starting app on port %s for your pleasure \n", portNum))
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

func run() error {

	// Reservation model stored in session
	gob.Register(models.Reservation{})

	//change to true when in production
	appCnf.InProduction = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	appCnf.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	appCnf.ErrorLog = errorLog

	//set up session
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appCnf.InProduction

	appCnf.Session = session

	tmplCache, err := render.CreateTemplateCache()
	if err != nil {
		fmt.Printf("Cannot create template cache, error %s \n", err)
		return err
		//fmt.Println(fmt.Sprintf("Error crating template configuration, error %s \n", err))
	}
	appCnf.TemplateCache = tmplCache
	appCnf.UseCache = false

	repo := handlers.NewRepo(&appCnf)
	handlers.NewHandlers(repo)

	render.NewTemplates(&appCnf)
	helpers.NewHelpers(&appCnf)

	return nil
}
