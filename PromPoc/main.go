package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "github.com/go-sql-driver/mysql"
)

var REQUEST_COUNT = promauto.NewCounter(prometheus.CounterOpts{
	Name: "go_app_requests_count",
	Help: "Total App HTTP Requests count.",
})
var REQUEST_INPROGRESS = promauto.NewGauge(prometheus.GaugeOpts{
	Name: "go_app_requests_inprogress",
	Help: "Number of application requests in progress",
})
var REQUEST_RESPOND_TIME = promauto.NewSummary(prometheus.SummaryOpts{
	Name: "go_app_response_latency_seconds",
	Help: "Response latency in seconds",
})
var fname []string
var lname []string
var mobile []string
var OPD []string
var issue []string
var date []string
var Time []string
var email []string
var sex []string

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", appoint)

	router.HandleFunc("/final", final)

	router.HandleFunc("/appointOrtho", appointOrtho)

	router.HandleFunc("/appointDiabet", appointDiabet)

	router.HandleFunc("/appointPedia", appointPedia)
	http.Handle("/", router)

	db, err := sql.Open("mysql", "root:vikas@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err)
	}

	query := "CREATE TABLE IF NOT EXISTS patients(PID int NOT NULL AUTO_INCREMENT,fname varchar(200),lname varchar(200),sex varchar(200),Mobile varchar(200),email varchar(200),OPD varchar(200),Issue varchar(200),Date varchar(200),Time varchar(200),PRIMARY KEY(PID))"
	create, err := db.Exec(query)
	if err != nil {
		panic(err)
	}
	fmt.Println(create)

	router.Path("/metrics").Handler(promhttp.Handler())
	http.ListenAndServe(":8000", nil)

}
func appoint(w http.ResponseWriter, r *http.Request) {

	REQUEST_COUNT.Inc()
	REQUEST_INPROGRESS.Inc()
	start_time := time.Now()
	time.Sleep(2 * time.Second)
	time_taken := time.Since(start_time)
	REQUEST_RESPOND_TIME.Observe(time_taken.Seconds())

	tmplt := template.Must(template.ParseFiles("templates/appoint.html"))

	tmplt.Execute(w, nil)

}
func final(w http.ResponseWriter, r *http.Request) {
	REQUEST_INPROGRESS.Dec()
	r.ParseForm()
	fname = r.Form["fname"]
	lname = r.Form["lname"]
	mobile = r.Form["mobile"]
	email = r.Form["email"]
	OPD = r.Form["OPD"]
	sex = r.Form["sex"]
	issue = r.Form["issue"]
	fmt.Println(fname, lname, mobile, OPD[0], issue)

	if check_user(mobile[0], email[0]) {
		var tmplt = template.Must(template.ParseFiles("templates/already.html"))
		tmplt.Execute(w, nil)
	} else {

		if OPD[0] == "Orthopedic" {
			var tmplt = template.Must(template.ParseFiles("templates/appointOrtho.html"))
			tmplt.Execute(w, nil)

		} else if OPD[0] == "Diabetes" {
			var tmplt = template.Must(template.ParseFiles("templates/appointDiabet.html"))
			tmplt.Execute(w, nil)

		} else {
			var tmplt = template.Must(template.ParseFiles("templates/appointPedia.html"))
			tmplt.Execute(w, nil)

		}

	}

}
func appointDiabet(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	Time = r.Form["time"]
	date = r.Form["date"]
	fmt.Println(Time, date)

	if check_appoint(Time[0], date[0]) {
		var tmplt = template.Must(template.ParseFiles("templates/chTmDt.html"))
		tmplt.Execute(w, nil)

	} else if add_user(fname[0], lname[0], sex[0], mobile[0], email[0], OPD[0], issue[0], Time[0], date[0]) {
		mail(email[0], fname[0], lname[0], sex[0], "Naga Vikas", OPD[0], Time[0], date[0])
		var tmplt = template.Must(template.ParseFiles("templates/index.html"))
		d := struct {
			First  string
			Last   string
			Doctor string
			Time   string
			Date   string
		}{
			First:  fname[0],
			Last:   lname[0],
			Doctor: "Naga Vikas",
			Time:   Time[0],
			Date:   date[0],
		}
		tmplt.Execute(w, d)
	}

}
func appointOrtho(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	Time = r.Form["time"]
	date = r.Form["date"]
	if check_appoint(Time[0], date[0]) {
		var tmplt = template.Must(template.ParseFiles("templates/chTmDt.html"))
		tmplt.Execute(w, nil)

	} else if add_user(fname[0], lname[0], sex[0], mobile[0], email[0], OPD[0], issue[0], Time[0], date[0]) {
		var tmplt = template.Must(template.ParseFiles("templates/index.html"))
		mail(email[0], fname[0], lname[0], sex[0], "Praveen Kumar", OPD[0], Time[0], date[0])
		d := struct {
			First  string
			Last   string
			Doctor string
			Time   string
			Date   string
		}{
			First:  fname[0],
			Last:   lname[0],
			Doctor: "Praveen Kumar",
			Time:   Time[0],
			Date:   date[0],
		}
		tmplt.Execute(w, d)
	}

}
func appointPedia(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	Time = r.Form["time"]
	date = r.Form["date"]
	if check_appoint(Time[0], date[0]) {
		var tmplt = template.Must(template.ParseFiles("templates/chTmDt.html"))
		tmplt.Execute(w, nil)

	} else if add_user(fname[0], lname[0], sex[0], mobile[0], email[0], OPD[0], issue[0], Time[0], date[0]) {
		mail(email[0], fname[0], lname[0], sex[0], "Khadar Basha", OPD[0], Time[0], date[0])
		var tmplt = template.Must(template.ParseFiles("templates/index.html"))
		d := struct {
			First  string
			Last   string
			Doctor string
			Time   string
			Date   string
		}{
			First:  fname[0],
			Last:   lname[0],
			Doctor: "Khadar Basha",
			Time:   Time[0],
			Date:   date[0],
		}
		tmplt.Execute(w, d)
	}

}

func check_user(mobile string, email string) bool {
	db, err := sql.Open("mysql", "root:vikas@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err)
	}
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT fname FROM patients WHERE Mobile='%s' AND email='%s')", (mobile), (email))
	row := db.QueryRow(query).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	fmt.Println(row)
	defer db.Close()
	return exists
}
func add_user(fname string, lname string, sex string, mobile string, email string, OPD string, issue string, time string, date string) bool {
	db, err := sql.Open("mysql", "root:vikas@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err)
	}
	add, err := db.Query("INSERT INTO patients(fname,lname,sex,Mobile,email,OPD,Issue,Date,Time) VALUES (?,?,?,?,?,?,?,?,?)", (fname), (lname), (sex), (mobile), (email), (OPD), (issue), (date), (time))
	if err != nil {
		panic(err)
	}
	fmt.Println(add)
	defer db.Close()
	return true
}
func check_appoint(time string, date string) bool {
	db, err := sql.Open("mysql", "root:vikas@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err)
	}
	var exists bool
	query := fmt.Sprintf("SELECT EXISTS(SELECT OPD FROM patients WHERE Time='%s' AND Date='%s')", (time), (date))
	row := db.QueryRow(query).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}
	fmt.Println(row)
	defer db.Close()
	return exists

}
func mail(email string, fname string, lname string, sex string, doctor string, OPD string, time string, date string) bool {

	to := []string{email}

	msg := "Your appointment is successful\n DETAILS:\nName:" + fname + " " + lname + "\nSex:" + sex + "\nDoctor:" + doctor + "\nSpecialist" + OPD + "\nslot:" + time + "\ndate:" + date + "\nAppointID:" + PID(email)
	message := []byte(msg)
	auth := smtp.PlainAuth("", "v26kas@gmail.com", "mirrigntwcikjyvl", "smtp.gmail.com")
	err := smtp.SendMail("smtp.gmail.com:587", auth, "v26kas@gmail.com", to, message)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Println("Email Sent Successfully!")
	return true
}
func PID(email string) string {
	db, err := sql.Open("mysql", "root:vikas@tcp(127.0.0.1:3306)/test")
	if err != nil {
		panic(err)
	}
	var PID string
	err = db.QueryRow("select PID from patients where email = ?", email).Scan(&PID)
	if err != nil {
		log.Fatal(err)
	}
	return PID

}
