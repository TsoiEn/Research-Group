package home

import (
	"database/sql"
	"html/template"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Global variable for DB connection
var db *sql.DB

// FOR ADMIN
// Handler to go into ADMIN CREDENTIALS PAGE
func adminCredHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("../FrontEnd/HomePage/AdminHomePage/AdminCredentials/AdminCredentials.html"))
	tmpl.Execute(w, nil)
}

// FOR STUDENT
// Handler to go into STUDENT CREDENTIALS PAGE
func stuCredHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("../FrontEnd/HomePage/StudentHomePage/StudentCredentials/StudentCred.html"))
	tmpl.Execute(w, nil)
}

func MainHome(database *sql.DB) *mux.Router {
	db = database // Assign the DB connection to the global variable

	// Set up routes (r)
	r := mux.NewRouter()

	// ROUTES
	r.HandleFunc("/admincredentialmanager", adminCredHandler).Methods("GET") // To ADMIN CREDENTIALS PAGE
	r.HandleFunc("/studentcredentials", stuCredHandler).Methods("GET")       // To STUDENT CREDENTIALS PAGE

	// STATIC FILES (Serve CSS, JS, images)
	fsAdminCred := http.FileServer(http.Dir("../FrontEnd/HomePage/AdminHomePage/AdminCredentials"))
	r.PathPrefix("/AdminCredentials/").Handler(http.StripPrefix("/AdminCredentials", fsAdminCred))

	fsStuCred := http.FileServer(http.Dir("../FrontEnd/HomePage/StudentHomePage/StudentCredentials"))
	r.PathPrefix("/StudentCredentials/").Handler(http.StripPrefix("/StudentCredentials", fsStuCred))

	return r
}
