//Handler package functions
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/t-Ikonen/bbbookingsystem/internal/config"
	"github.com/t-Ikonen/bbbookingsystem/internal/driver"
	"github.com/t-Ikonen/bbbookingsystem/internal/forms"
	"github.com/t-Ikonen/bbbookingsystem/internal/helpers"
	"github.com/t-Ikonen/bbbookingsystem/internal/models"
	"github.com/t-Ikonen/bbbookingsystem/internal/render"
	"github.com/t-Ikonen/bbbookingsystem/internal/repository"
	"github.com/t-Ikonen/bbbookingsystem/internal/repository/dbrepo"
)

// Repo used by handlers
var Repo *Repository

//Repository is repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

//NewRepo a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

//NewTestRepo a new testing repository
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

//NewHandlers sets the repository for the Handlers
func NewHandlers(r *Repository) {
	Repo = r
}

//Home page function hadles Home page
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "home.page.tmpl.html", &models.TemplateData{}, r)
}

// About func handles About page
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "about.page.tmpl.html", &models.TemplateData{}, r)
}

//Booking to render Booking page
func (m *Repository) Booking(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "booking.page.tmpl.html", &models.TemplateData{}, r)
}

//PostBooking to post Booking page data
func (m *Repository) PostBooking(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	end := r.Form.Get("endDate")
	start := r.Form.Get("startDate")

	layout := "01-02-2006"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	for _, i := range rooms {
		m.App.InfoLog.Println("ROOM:", i.ID, i.RoomName)
	}

	if len(rooms) == 0 {
		m.App.InfoLog.Println("No rooms")
		m.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/booking", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, "chooseroom.page.tmpl.html", &models.TemplateData{
		Data: data,
	}, r)

	//w.Write([]byte(fmt.Sprintf("Posting start: %s and end is %s", start, end)))

}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `jsnon:"message"`
	RoomId    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

//BookingJSON to request availability JSON format
func (m *Repository) BookingJSON(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		// can't parse form, so return appropriate json
		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error",
		}

		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}
	sd := r.Form.Get("start")
	ed := r.Form.Get("end")
	layout := "01-02-2006"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		//log.Println(err)
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(layout, ed)
	if err != nil {
		//log.Println(err)
		helpers.ServerError(w, err)
		return
	}
	roomId, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		//log.Println(err)
		helpers.ServerError(w, err)
		return
	}
	available, err := m.DB.SearchAvailabilityByDatesByRoomId(startDate, endDate, roomId)
	if err != nil {
		//log.Println(err)
		helpers.ServerError(w, err)
		return
	}
	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomId:    strconv.Itoa(roomId),
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		//log.Println(err)
		helpers.ServerError(w, err)
		return
	}
	//log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

//Reservation to render Reservation page
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "cannot get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := m.DB.GetRoomNameById(res.RoomId)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "Cannot find room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	res.Room.RoomName = room.RoomName

	m.App.Session.Put(r.Context(), "reservation", res)

	//log.Println("room name: ", res.Room.RoomName)
	//log.Println("Render reservation room id: ", res.RoomId)
	sd := res.StartDate.Format("01-02-2006")
	ed := res.EndDate.Format("01-02-2006")

	stringmap := make(map[string]string)
	stringmap["start_date"] = sd
	stringmap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res
	render.Template(w, "reservation.page.tmpl.html", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringmap,
	}, r)
}

//PostReservation handel posting Reservetation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation data from session"))
		return
	}

	err := r.ParseForm()

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot parse form")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Phone = r.Form.Get("phone")
	reservation.Email = r.Form.Get("email")

	form := forms.New(r.PostForm)

	// form.Has("first_name", r)
	form.Required("first_name", "last_name", "email")

	form.MinLenght("first_name", 3)
	form.ValidEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, "reservation.page.tmpl.html", &models.TemplateData{
			Form: form,
			Data: data,
		}, r)
		return
	}

	newReservationId, err := m.DB.InsertReservation(reservation)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot insert reservation to DB")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//m.App.Session.Put(r.Context(), "reservation", reservation)

	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndSate:       reservation.EndDate,
		RoomId:        reservation.RoomId,
		ReservationId: newReservationId,
		RestrictionId: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot insert room restriction to DB")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//subject := "Reservation confirmation" + strconv.(reservation.StartDate) + "-" + reservation.EndDate
	//send email notification - first to guest
	customerMessage := fmt.Sprintf(`
		<strong>Reservation confirmation</strong><br><br>
		Dear %s %s, <br><hr>
		This is to confirm your reservation for %s from %s to %s.
		`, reservation.FirstName, reservation.LastName, reservation.Room.RoomName, reservation.StartDate.Format("02-01-2006"), reservation.EndDate.Format("02-01-2006"))

	customerMessage2 := fmt.Sprintln(`
	<br><br><br>Welcome to be reborn again!<br>
	<hr>
	Contact: black.lodge@outlook.com | wwww.blacklodge.xyz
	`)

	msg := models.MailData{
		To:       reservation.Email,
		From:     "ed.glen@blacklodge.xyz",
		Subject:  "Reservation confirmation",
		Message:  customerMessage + customerMessage2,
		Template: "basic.html",
	}
	m.App.MailChan <- msg

	//send email notification - to guest owner
	ownerMessage := fmt.Sprintf(`
	<strong>Reservation notification</strong><br>
	Room reservation notification <br><hr><br><br>
	First name: %s <br>
	Last name: %s<br>
	Room: %s<br>
	Start date: %s<br>
	End date: %s<br>
	`, reservation.FirstName, reservation.LastName, reservation.Room.RoomName, reservation.StartDate.Format("02-01-2006"), reservation.EndDate.Format("02-01-2006"))

	msg2 := models.MailData{
		To:      reservation.Email,
		From:    "ed.glen@blacklodge.xyz",
		Subject: "Reservation notification",
		Message: ownerMessage,
	}
	m.App.MailChan <- msg2

	m.App.Session.Put(r.Context(), "reservation", reservation)
	http.Redirect(w, r, "/reservationsummary", http.StatusSeeOther)

}

//ChooseRoom to render Choose Room page that lists availabe room
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	// roomId, err := strconv.Atoi(chi.URLParam(r, "id"))
	// if err != nil {
	// 	m.App.Session.Put(r.Context(), "error", "cannot convert room id to integer")
	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }

	// changed to this, so we can test it more easily
	// split the URL up by /, and grab the 3rd element
	exploded := strings.Split(r.RequestURI, "/")
	roomId, err := strconv.Atoi(exploded[2])
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	res.RoomId = roomId

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/reservation", http.StatusSeeOther)

}

//BookRoom takes URL parameters, makes reservation into session, redirects to make reservation page
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {

	roomId, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")
	//log.Println("room ID from URL:", roomId)
	// log.Println("start from URL:", sd)
	// log.Println("end from URL:", ed)

	layout := "01-02-2006"
	startDate, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "cannot parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var res models.Reservation
	//log.Println("getting room name by ID->", roomId, " <-")
	room, err := m.DB.GetRoomNameById(roomId)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	res.Room.RoomName = room.RoomName
	res.RoomId = roomId
	res.StartDate = startDate
	res.EndDate = endDate

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/reservation", http.StatusSeeOther)
}

//Contact to render contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "contact.page.tmpl.html", &models.TemplateData{}, r)
}

//Northernlights to render northernlights page
func (m *Repository) Northernlights(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "northernlights.page.tmpl.html", &models.TemplateData{}, r)
}

//Frostsuite to render Frostsuite page
func (m *Repository) Frostsuite(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "frostsuite.page.tmpl.html", &models.TemplateData{}, r)
}

//Snowsuite renders Snowsuite page
func (m *Repository) Snowsuite(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "snowsuite.page.tmpl.html", &models.TemplateData{}, r)
}

//Reservationsummary renders the reservation summary page
func (m *Repository) Reservationsummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	m.App.Session.Remove(r.Context(), "reservation")
	data := make(map[string]interface{})
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("01-02-2006")
	ed := reservation.EndDate.Format("01-02-2006")

	stringmap := make(map[string]string)
	stringmap["start_date"] = sd
	stringmap["end_date"] = ed

	render.Template(w, "reservationsummary.page.tmpl.html", &models.TemplateData{
		Data:      data,
		StringMap: stringmap,
	}, r)
}

//ShowLogin renders the Login page
func (m *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "login.page.tmpl.html", &models.TemplateData{
		Form: forms.New(nil),
	}, r)
}

//PostShowLogin posts the user login data and validates it
func (m *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	//log.Println("PostShowLogin works")
	_ = m.App.Session.RenewToken(r.Context())
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.ValidEmail("email")
	if !form.Valid() {
		render.Template(w, "login.page.tmpl.html", &models.TemplateData{
			Form: form,
		}, r)
		return
	}

	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "You are succesfully logged in.")
	http.Redirect(w, r, "/", http.StatusSeeOther)

}

//ShowLogout logs out user
func (m *Repository) ShowLogout(w http.ResponseWriter, r *http.Request) {
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

//AdminDashboard show dashboard page in admin tool
func (m *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "admindashboard.page.tmpl.html", &models.TemplateData{}, r)

}

//AdminNewReservations show new unprocessed reservations in admin tool
func (m *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {

	reservations, err := m.DB.AllNewReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["reservations"] = reservations
	render.Template(w, "adminnewreservations.page.tmpl.html", &models.TemplateData{
		Data: data,
	}, r)

}

//AdminAllReservations show all reservations in admin tool
func (m *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := m.DB.AllReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]interface{})
	data["reservations"] = reservations
	render.Template(w, "adminallreservations.page.tmpl.html", &models.TemplateData{
		Data: data,
	}, r)

}

//AdminCalendar show statistics for admin only
func (m *Repository) AdminCalendar(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "admincalendar.page.tmpl.html", &models.TemplateData{}, r)

}

//AdminShowReservation shows one reservation in admin tool for processing
func (m *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request) {
	explode := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(explode[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	src := explode[3]
	stringMap := make(map[string]string)
	stringMap["src"] = src

	//get reservaton from DB
	res, err := m.DB.GetReservationById(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["reservation"] = res
	//fmt.Println(data)
	render.Template(w, "adminshowreservation.page.tmpl.html", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		Form:      forms.New(nil),
	}, r)
}

//Save edited reservation in admin mode
func (m *Repository) AdminPostReservation(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	explode := strings.Split(r.RequestURI, "/")
	id, err := strconv.Atoi(explode[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	src := explode[3]
	stringMap := make(map[string]string)
	stringMap["src"] = src

	res, err := m.DB.GetReservationById(id)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.Phone = r.Form.Get("phone")
	//fmt.Println("update db alkaa")
	err = m.DB.UpdateReservation(res)
	//.Println("update db tehty")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	m.App.Session.Put(r.Context(), "flash", "Changes saved.")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}

//AdminStatistics show statistics for admin only
func (m *Repository) AdminStatistics(w http.ResponseWriter, r *http.Request) {
	render.Template(w, "statistics.page.tmpl.html", &models.TemplateData{}, r)

}

//AdminProcessReservation marks reservation as processed
func (m *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")
	_ = m.DB.UpdatePrcessed(id, 1)
	m.App.Session.Put(r.Context(), "flash", "Reservation marked as processed.")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)

}

//delete-reservation
//AdminDelteReservation deletes a reservation in admin mode
func (m *Repository) AdminDelteReservation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")
	_ = m.DB.DeleteReservation(id)
	m.App.Session.Put(r.Context(), "flash", "Reservation deleted")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)

}
