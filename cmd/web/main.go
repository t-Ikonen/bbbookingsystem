package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/t-Ikonen/bbbookingsystem/internal/config"
	"github.com/t-Ikonen/bbbookingsystem/internal/handlers"
	"github.com/t-Ikonen/bbbookingsystem/internal/render"
)

const portNum = ":8080"

var appCnf config.AppConfig
var session *scs.SessionManager

//Main of HelloWeb app
func main() {

	//change to true when in production
	appCnf.InProduction = false

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = appCnf.InProduction

	appCnf.Session = session

	tmplCache, err := render.CreateTemplateCache()
	if err != nil {
		fmt.Printf("Error crating template configuration, error %s \n", err)
		//fmt.Println(fmt.Sprintf("Error crating template configuration, error %s \n", err))
	}
	appCnf.TemplateCache = tmplCache
	appCnf.UseCache = false

	repo := handlers.NewRepo(&appCnf)
	handlers.NewHandlers(repo)

	render.NewTemplates(&appCnf)

	srv := &http.Server{
		Addr:    portNum,
		Handler: Routes(&appCnf),
	}
	fmt.Printf("Starting app on port %s for your pleasure \n", portNum)
	//fmt.Println(fmt.Sprintf("Starting app on port %s for your pleasure \n", portNum))
	err = srv.ListenAndServe()
	log.Fatal(err)

}
