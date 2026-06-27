package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

type Record struct {
	Date			string
	MountainName	string
	Weather			string
	Duration		string
	Memo			string
}

var db *sql.DB
var tmpl = template.Must(template.ParseFiles("templates/index.html"))

func main() {
	dsn := "root:Yoshimasa22@tcp(127.0.0.1:3306)/mountain_db"
	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatal("Cannot connect to DataBase: ", err)
	}

	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/add", handleAddRecord)

	fmt.Println("Running server: http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}


	rows, err := db.Query("SELECT date, mountain_name, weather, duration, memo FROM hiking_records ORDER BY date DESC")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var records []Record
	for rows.Next() {
		var rec Record
		if err := rows.Scan(&rec.Date, &rec.MountainName, &rec.Weather, &rec.Duration, &rec.Memo); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		records = append(records, rec)
	}

	tmpl.Execute(w, records)
}

func handleAddRecord(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// get data from the form
	date := r.FormValue("date")
	mountainName := r.FormValue("mountain_name")
	weather := r.FormValue("weather")
	duration := r.FormValue("duration")
	memo := r.FormValue("memo")

	if date == "" || mountainName == "" {
		http.Error(w, "Must contain Date and Mountain Name", http.StatusBadRequest)
		return
	}

	query := `
		INSERT INTO hiking_records (date, mountain_name, weather, duration, memo)
		VALUES (?, ?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE mountain_name=?, weather=?, duration=?, memo=?`

	_, err := db.Exec(query, date, mountainName, weather, duration, memo, mountainName, weather, duration, memo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
