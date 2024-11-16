package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

var db *sql.DB

// Initialize DB connection
func initDB() {
	var err error
	db, err = sql.Open("mysql", "RaffyEL:WeMakingDBS1@E@tcp(127.0.0.1:3306)/stucredstorage")
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MySQL!")
}

type Account struct {
	Username string
	Password string
}

// Handler to go into next page (TEST (will delete))
func successHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("success.html"))
	tmpl.Execute(w, nil)
}

func main() {
	// Initialize the database
	initDB()
	defer db.Close()

	// Set up routes (r)
	r := mux.NewRouter()

	// ROUTES
	r.HandleFunc("/adminlogin", adminHandler).Methods("GET", "POST")     // Admin login
	r.HandleFunc("/alumnilogin", alumniHandler).Methods("GET", "POST")   // Alumni login
	r.HandleFunc("/studentlogin", studentHandler).Methods("GET", "POST") // Student login
	r.HandleFunc("/success", successHandler).Methods("GET")              // Success page

	// STATIC FILES (Serve CSS, JS, images)
	fsAdmin := http.FileServer(http.Dir("../../FrontEnd/LoginPage/adminlog"))
	r.PathPrefix("/adminlog/").Handler(http.StripPrefix("/adminlog", fsAdmin))

	fsAlumni := http.FileServer(http.Dir("../../FrontEnd/LoginPage/alumnilog"))
	r.PathPrefix("/alumnilog/").Handler(http.StripPrefix("/alumnilog", fsAlumni))

	fsStud := http.FileServer(http.Dir("../../FrontEnd/LoginPage"))
	r.PathPrefix("/").Handler(http.StripPrefix("/", fsStud))

	// Start the local server
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
