//Handler package functions
package handlers

import (
	"fmt"
	"net/http"

	"github.com/t-Ikonen/bbbookingsystem/pkg/config"
	"github.com/t-Ikonen/bbbookingsystem/pkg/models"
	"github.com/t-Ikonen/bbbookingsystem/pkg/render"
)

// Repo used by handlers
var Repo *Repository

//Repository is repository type
type Repository struct {
	App *config.AppConfig
}

//NewRepo a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

//NewHandlers sets the repository for the Handlers
func NewHandlers(r *Repository) {
	Repo = r
}

//Home page function hadles Home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	ipAddr := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", ipAddr)

	fmt.Println("IP on home sivulla ", ipAddr)
	//send data to template
	render.RenderTemplate(w, "home.page.tmpl.html", &models.TemplateData{})
}

// About func handles About page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello again"

	remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIp
	fmt.Println("IP on About sivulla ", remoteIp)

	render.RenderTemplate(w, "about.page.tmpl.html", &models.TemplateData{
		StringMap: stringMap,
	})
}
