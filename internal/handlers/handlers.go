package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/AlessioPani/go-booking/internal/config"
	"github.com/AlessioPani/go-booking/internal/driver"
	"github.com/AlessioPani/go-booking/internal/forms"
	"github.com/AlessioPani/go-booking/internal/helpers"
	"github.com/AlessioPani/go-booking/internal/models"
	"github.com/AlessioPani/go-booking/internal/renders"
	"github.com/AlessioPani/go-booking/internal/repository"
	"github.com/AlessioPani/go-booking/internal/repository/dbrepo"
	"github.com/go-chi/chi/v5"
)

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// Repo is the repository used by handlers
var Repo *Repository

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

// NewTestRepo creates a new repository for testing purpose
func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestRepo(a),
	}
}

// NewHandlers sets the repository for handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the homepage handler.
func (pr *Repository) Home(w http.ResponseWriter, r *http.Request) {
	renders.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler.
func (pr *Repository) About(w http.ResponseWriter, r *http.Request) {
	renders.Template(w, r, "about.page.tmpl", &models.TemplateData{})
}

// Reservation renders the make a reservation page and displays form
func (pr *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	res, ok := pr.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		pr.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := pr.DB.GetRoomById(res.RoomId)
	if err != nil {
		pr.App.Session.Put(r.Context(), "error", "Can't find room")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName
	pr.App.Session.Put(r.Context(), "reservation", res)

	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]any)
	data["reservation"] = res

	renders.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// PostReservation handles the posting of a reservation form
func (pr *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		pr.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")

	// 2020-01-01 -- 01/02 03:04:05PM '06 -0700
	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		pr.App.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		pr.App.Session.Put(r.Context(), "error", "can't get parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		pr.App.Session.Put(r.Context(), "error", "invalid data!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res, ok := pr.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		pr.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Phone:     r.Form.Get("phone"),
		Email:     r.Form.Get("email"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomId:    roomID,
		Room:      res.Room,
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]any)
		data["reservation"] = reservation
		http.Error(w, "my own error message", http.StatusSeeOther)
		renders.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	newReservationID, err := pr.DB.InsertReservation(reservation)
	if err != nil {
		pr.App.Session.Put(r.Context(), "error", "can't insert reservation into database!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	restriction := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomId:        roomID,
		ReservationId: newReservationID,
		RestrictionId: 1,
	}

	err = pr.DB.InsertRoomRestriction(restriction)
	if err != nil {
		pr.App.Session.Put(r.Context(), "error", "can't insert room restriction!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// send mail notification - first to guest
	htmlMessage := fmt.Sprintf(
		`<strong>Reservation Confirmation</strong><br><br>
		Dear %s, <br>
		this is a confirmation of your reservation from %s to %s .<br><br>

		Looking forward to see you soon<br><br>
		Best regards,<br>
		Admin`,
		reservation.FirstName, reservation.StartDate.Format("2006-01-02"),
		reservation.EndDate.Format("2006-01-02"))

	msg := models.MailData{
		To:       reservation.Email,
		From:     "reservation@me.com",
		Subject:  "Reservation Confirmation",
		Content:  htmlMessage,
		Template: "basic.html",
	}

	pr.App.MailChan <- msg

	// send mail notification - second to owner
	htmlMessage = fmt.Sprintf(
		`<strong>Reservation Notification</strong><br><br>
		Dear Admin, <br>
		there is a new reservation from Mr./Mrs. %s %s, from %s to %s. <br><br><br>

		Kind regards,<br>
		Admin`,
		reservation.FirstName, reservation.LastName, reservation.StartDate.Format("2006-01-02"),
		reservation.EndDate.Format("2006-01-02"))

	msgToAdmin := models.MailData{
		To:      "me@me.com",
		From:    "me@me.com",
		Subject: "Reservation Notice",
		Content: htmlMessage,
	}

	pr.App.MailChan <- msgToAdmin

	pr.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Generals renders the room page
func (pr *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	renders.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

// Majors renders the room page
func (pr *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	renders.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

// Availability renders the search availability page
func (pr *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	renders.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability renders the search availability page
func (pr *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		pr.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	start := r.Form.Get("start")
	end := r.Form.Get("end")

	// 2020-01-01 --- Format using this 01/02 03:04:05PM '06 -0700 as a reference
	layout := "2006-01-02"
	startDate, err := time.Parse(layout, start)
	if err != nil {
		pr.App.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		pr.App.Session.Put(r.Context(), "error", "can't parse end date")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	rooms, err := pr.DB.SearchAvailabilityForAllRooms(startDate, endDate)
	if err != nil {
		pr.App.Session.Put(r.Context(), "error", "can't get availability for rooms")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	for _, i := range rooms {
		pr.App.InfoLog.Println("Room: ", i.ID, i.RoomName)
	}

	if len(rooms) == 0 {
		// no availability
		pr.App.Session.Put(r.Context(), "error", "No availability")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]any)
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	pr.App.Session.Put(r.Context(), "reservation", res)

	renders.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})

}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// AvailabilityJSON handlers requests for availability and send JSON response
func (pr *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Internal server error",
		}

		out, _ := json.MarshalIndent(resp, "", "  ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	sd := r.Form.Get("start")
	ed := r.Form.Get("end")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)
	roomId, _ := strconv.Atoi(r.Form.Get("room_id"))

	available, err := pr.DB.SearchAvailabilityByDatesByRoomId(startDate, endDate, roomId)
	if err != nil {
		resp := jsonResponse{
			OK:      false,
			Message: "Error connecting to the database",
		}

		out, _ := json.MarshalIndent(resp, "", "  ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: sd,
		EndDate:   ed,
		RoomID:    strconv.Itoa(roomId),
	}

	// removed error checking, resp is created manually
	out, _ := json.MarshalIndent(resp, "", " ")

	w.Header().Set("Content-type", "application/json")
	w.Write(out)
}

// Contact renders the contact page
func (pr *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	renders.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// ReservationSummary displays the reservation summary page
func (pr *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := pr.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		pr.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	pr.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]any)
	data["reservation"] = reservation

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")
	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	renders.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

// ChooseRoom displays a list of available rooms
func (pr *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	roomID, err := strconv.Atoi(exploded[2])
	if err != nil {
		pr.App.Session.Put(r.Context(), "error", "missing url parameter")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res, ok := pr.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		pr.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.RoomId = roomID

	pr.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

// BookRoom takes url parameters, builds a session variable and takes user to the make reservation page
func (pr *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomId, _ := strconv.Atoi(r.URL.Query().Get("id"))
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "2006-01-02"
	startDate, _ := time.Parse(layout, sd)
	endDate, _ := time.Parse(layout, ed)

	// create a reservation
	var res models.Reservation

	room, err := pr.DB.GetRoomById(roomId)
	if err != nil {
		pr.App.Session.Put(r.Context(), "error", "can't get room by id")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.StartDate = startDate
	res.EndDate = endDate
	res.RoomId = roomId
	res.Room.RoomName = room.RoomName

	pr.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

// ShowLogin shows login page
func (pr *Repository) ShowLogin(w http.ResponseWriter, r *http.Request) {
	renders.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// PostShowLogin handles logins of users
func (pr *Repository) PostShowLogin(w http.ResponseWriter, r *http.Request) {
	_ = pr.App.Session.RenewToken(r.Context())
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	if !form.Valid() {
		renders.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")
	id, _, err := pr.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)
		pr.App.Session.Put(r.Context(), "error", "Invalid login credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	pr.App.Session.Put(r.Context(), "user_id", id)
	pr.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout logs a user out
func (pr *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	_ = pr.App.Session.Destroy(r.Context())
	pr.App.Session.RenewToken(r.Context())
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// AdminDashboard shows an admin page - required authentication
func (pr *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	renders.Template(w, r, "admin-dashboard.page.tmpl", &models.TemplateData{})
}

// AdminNewReservations shows an admin page with all new reservations
func (pr *Repository) AdminNewReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := pr.DB.AllNewReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]any)
	data["reservations"] = reservations
	renders.Template(w, r, "admin-reservations-new.page.tmpl", &models.TemplateData{Data: data})
}

// AdminAllReservations shows an admin page with all reservations
func (pr *Repository) AdminAllReservations(w http.ResponseWriter, r *http.Request) {
	reservations, err := pr.DB.AllReservations()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	data := make(map[string]any)
	data["reservations"] = reservations
	renders.Template(w, r, "admin-reservations-all.page.tmpl", &models.TemplateData{Data: data})
}

// AdminCalendarReservations displays a reservations calendar
func (pr *Repository) AdminCalendarReservations(w http.ResponseWriter, r *http.Request) {
	// assume that there is no month or year specified
	now := time.Now()

	if r.URL.Query().Get("y") != "" {
		year, _ := strconv.Atoi(r.URL.Query().Get("y"))
		month, _ := strconv.Atoi(r.URL.Query().Get("m"))
		now = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	}

	data := make(map[string]any)
	data["now"] = now

	next := now.AddDate(0, 1, 0)
	last := now.AddDate(0, -1, 0)

	nextMonth := next.Format("01")
	nextMonthYear := next.Format("2006")

	lastMonth := last.Format("01")
	lastMonthYear := last.Format("2006")

	stringMap := make(map[string]string)
	stringMap["next_month"] = nextMonth
	stringMap["next_month_year"] = nextMonthYear
	stringMap["last_month"] = lastMonth
	stringMap["last_month_year"] = lastMonthYear

	stringMap["this_month"] = now.Format("01")
	stringMap["this_month_year"] = now.Format("2006")

	// get first and last date of a month
	currentYear, currentMonth, _ := now.Date()
	currentLocation := now.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 0, 0, 0, 0, currentLocation)
	lastofMonth := firstOfMonth.AddDate(0, 1, -1)

	intMap := make(map[string]int)
	intMap["days_in_month"] = lastofMonth.Day()

	rooms, err := pr.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data["rooms"] = rooms

	for _, x := range rooms {
		// create maps
		reservationMap := make(map[string]int)
		blockMap := make(map[string]int)

		for d := firstOfMonth; !d.After(lastofMonth); d = d.AddDate(0, 0, 1) {
			reservationMap[d.Format("2006-01-2")] = 0
			blockMap[d.Format("2006-01-2")] = 0
		}

		// get all restriction for the current room
		restrictions, err := pr.DB.GetRestrictionsForRoomByDate(x.ID, firstOfMonth, lastofMonth)
		if err != nil {
			helpers.ServerError(w, err)
			return
		}

		for _, y := range restrictions {
			if y.ReservationId > 0 {
				// it's a reservation
				for d := y.StartDate; !d.After(y.EndDate); d = d.AddDate(0, 0, 1) {
					reservationMap[d.Format("2006-01-2")] = y.ReservationId
				}
			} else {
				// it's a block
				for d := y.StartDate; !d.After(y.EndDate); d = d.AddDate(0, 0, 1) {
					blockMap[d.Format("2006-01-2")] = y.ID
				}
			}
		}

		data[fmt.Sprintf("reservation_map_%d", x.ID)] = reservationMap
		data[fmt.Sprintf("block_map_%d", x.ID)] = blockMap

		pr.App.Session.Put(r.Context(), fmt.Sprintf("block_map_%d", x.ID), blockMap)
	}

	renders.Template(w, r, "admin-reservation-calendar.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		IntMap:    intMap,
	})
}

// AdminPostCalendarReservations handles post of reservation calendar
func (pr *Repository) AdminPostCalendarReservations(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	year, _ := strconv.Atoi(r.Form.Get("y"))
	month, _ := strconv.Atoi(r.Form.Get("m"))

	// process blocks
	rooms, err := pr.DB.AllRooms()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	form := forms.New(r.PostForm)
	for _, x := range rooms {
		// Get the block map from session. Loop through the entire map, if we have an entry in the map
		// that does not exist in our posted data, and if restriction id > 0, then it is a block we need
		// to remove.
		curMap := pr.App.Session.Get(r.Context(), fmt.Sprintf("block_map_%d", x.ID)).(map[string]int)
		for name, value := range curMap {
			// ok will be false if the value is not in the map
			if val, ok := curMap[name]; ok {
				// only pay attention to values > 0, and that are not in the form post
				// the rest are just placeholders for days without blocks
				if val > 0 {
					if !form.Has(fmt.Sprintf("remove_block_%d_%s", x.ID, name)) {
						//delete restriction by id
						err := pr.DB.DeleteBlockById(value)
						if err != nil {
							log.Println(err)
						}
					}
				}
			}
		}
	}

	// Handle new blocks
	for name, _ := range r.PostForm {
		if strings.HasPrefix(name, "add_block") {
			exploded := strings.Split(name, "_")
			roomId, _ := strconv.Atoi(exploded[2])
			date, err := time.Parse("2006-01-2", exploded[3])

			// insert a new block
			err = pr.DB.AddBlockForRoom(roomId, date)
			if err != nil {
				log.Println(err)
			}
		}
	}

	pr.App.Session.Put(r.Context(), "flash", "Changes saved")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-cal?y=%d&m=%d", year, month), http.StatusSeeOther)
}

// AdminShowReservation displays in a detailed view a single reservation
func (pr *Repository) AdminShowReservation(w http.ResponseWriter, r *http.Request) {
	exploded := strings.Split(r.RequestURI, "/")
	roomID, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := exploded[3]
	stringMap := make(map[string]string)
	stringMap["src"] = src

	// get reservation from database
	res, err := pr.DB.GetReservationById(roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]any)
	data["reservation"] = res

	renders.Template(w, r, "admin-reservation-show.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
		Data:      data,
		Form:      forms.New(nil),
	})
}

// AdminPostShowReservation displays in a detailed view a single reservation
func (pr *Repository) AdminPostShowReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	exploded := strings.Split(r.RequestURI, "/")
	roomID, err := strconv.Atoi(exploded[4])
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	src := exploded[3]

	// get reservation from database
	res, err := pr.DB.GetReservationById(roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.FirstName = r.Form.Get("first_name")
	res.LastName = r.Form.Get("last_name")
	res.Email = r.Form.Get("email")
	res.Phone = r.Form.Get("phone")

	err = pr.DB.UpdateReservation(res)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	pr.App.Session.Put(r.Context(), "flash", "Changes saved")
	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}

// AdminProcessReservation marks a reservation as processed
func (pr *Repository) AdminProcessReservation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")

	_ = pr.DB.UpdatedProcessedForReservation(id, 1)

	pr.App.Session.Put(r.Context(), "flash", "Reservation marked as processed")

	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)

}

// AdminDeleteReservation deletes a reservation
func (pr *Repository) AdminDeleteReservation(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	src := chi.URLParam(r, "src")

	_ = pr.DB.DeleteReservation(id)

	pr.App.Session.Put(r.Context(), "flash", "Reservation deleted")

	http.Redirect(w, r, fmt.Sprintf("/admin/reservations-%s", src), http.StatusSeeOther)
}
