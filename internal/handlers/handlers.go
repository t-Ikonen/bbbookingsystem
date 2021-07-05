//Handler package functions
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/t-Ikonen/bbbookingsystem/internal/config"
	"github.com/t-Ikonen/bbbookingsystem/internal/forms"
	"github.com/t-Ikonen/bbbookingsystem/internal/models"
	"github.com/t-Ikonen/bbbookingsystem/internal/render"
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
	render.RenderTemplate(w, "home.page.tmpl.html", &models.TemplateData{}, r)
}

// About func handles About page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	// stringMap := make(map[string]string)
	// stringMap["test"] = "Hello again"

	// remoteIp := m.App.Session.GetString(r.Context(), "remote_ip")
	// stringMap["remote_ip"] = remoteIp
	// fmt.Println("IP on About sivulla ", remoteIp)

	render.RenderTemplate(w, "about.page.tmpl.html", &models.TemplateData{
		//StringMap: stringMap,
	}, r)
}

//Booking to render Booking page
func (m *Repository) Booking(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "booking.page.tmpl.html", &models.TemplateData{}, r)
}

//PostBooking to post Booking page data
func (m *Repository) PostBooking(w http.ResponseWriter, r *http.Request) {
	end := r.Form.Get("endDate")
	start := r.Form.Get("startDate")
	w.Write([]byte(fmt.Sprintf("Posting start: %s and end is %s", start, end)))
	//w.Write([]byte("postiiiiing"))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `jsnon:"message"`
}

//BookingJSON to request availability JSON format
func (m *Repository) BookingJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Not Available",
	}
	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		log.Println(err)
	}
	log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

//Reservation to render Reservation page
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation

	data := make(map[string]interface{})
	data["reservation"] = emptyReservation
	render.RenderTemplate(w, "reservation.page.tmpl.html", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	}, r)
}

//PostReservation handel posting Reservetation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}
	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}
	form := forms.New(r.PostForm)

	// form.Has("first_name", r)
	form.Required("first_name", "last_name", "email")

	form.MinLenght("first_name", 3, r)
	form.ValidEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.RenderTemplate(w, "reservation.page.tmpl.html", &models.TemplateData{
			Form: form,
			Data: data,
		}, r)
		return
	}

}

//Contact to render contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "contact.page.tmpl.html", &models.TemplateData{}, r)
}

//Northernlights to render northernlights page
func (m *Repository) Northernlights(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "northernlights.page.tmpl.html", &models.TemplateData{}, r)
}

//Frostsuite to render Frostsuite page
func (m *Repository) Frostsuite(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "frostsuite.page.tmpl.html", &models.TemplateData{}, r)
}

//Snowsuite renders Snowsuite page
func (m *Repository) Snowsuite(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "snowsuite.page.tmpl.html", &models.TemplateData{}, r)
}
