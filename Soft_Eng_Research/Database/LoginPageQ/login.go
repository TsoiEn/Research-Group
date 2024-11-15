package main

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
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

// Handler to fetch account information
func loginHandler(w http.ResponseWriter, r *http.Request) {
	var errorMessage string

	if r.Method == http.MethodPost {
		// Parse form data
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Check information in the database
		var storedPassword, accountID string
		err := db.QueryRow("SELECT password, accountID FROM accounts WHERE username = ?", username).Scan(&storedPassword, &accountID)
		if err != nil {
			if err == sql.ErrNoRows {
				errorMessage = "Invalid username or password"
			} else {
				errorMessage = "Server error. Please try again later."
			}
		} else {
			// Hash the input password for accurate comparison
			hash := sha256.Sum256([]byte(password))
			hashedPassword := hex.EncodeToString(hash[:])

			// Compare the hashed input password and the hashed stored password
			if hashedPassword != storedPassword {
				errorMessage = "Invalid username or password."
			} else {
				// Compare if the accountID starts with "3", indicating a student account
				if accountID[:1] != "3" {
					errorMessage = "This is not student account."
				} else {
					// Render next page if it's a valid student account
					http.Redirect(w, r, "/success", http.StatusSeeOther)
					return
				}
			}
		}

		// If there is an error, render the login page with the error message
		tmpl := template.Must(template.ParseFiles("../../FrontEnd/LoginPage/login.html"))
		tmpl.Execute(w, struct {
			ErrorMessage string
			Username     string
			Password     string
		}{
			ErrorMessage: errorMessage,
			Username:     username,
			Password:     "", //Clear the password field
		})
		return
	}

	// Initial form rendering
	tmpl := template.Must(template.ParseFiles("../../FrontEnd/LoginPage/login.html"))
	tmpl.Execute(w, nil)
}

// Handler to go into next page
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

	// Define routes for login and next page
	r.HandleFunc("/login", loginHandler).Methods("GET", "POST")
	r.HandleFunc("/success", successHandler).Methods("GET")

	// Serve static files (CSS, JS, images)
	fs := http.FileServer(http.Dir("../../FrontEnd/LoginPage"))
	r.PathPrefix("/").Handler(http.StripPrefix("/", fs))

	// Start the local server
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
