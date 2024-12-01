package login

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"html/template"
	"net/http"
)

// Handler to fetch ALUMNI account information
func alumniHandler(w http.ResponseWriter, r *http.Request) {
	var errorMessage string

	if r.Method == http.MethodPost {
		// Parse form data
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Query the database for alumni credentials
		var storedPassword, accountID string
		err := db.QueryRow("SELECT password, accountID FROM accounts WHERE username = ?", username).Scan(&storedPassword, &accountID)
		if err != nil {
			if err == sql.ErrNoRows {
				errorMessage = "Invalid username or password"
			} else {
				errorMessage = "Server error. Please try again later."
			}
		} else {
			// Hash the input password for comparison
			hash := sha256.Sum256([]byte(password))
			hashedPassword := hex.EncodeToString(hash[:])

			// Compare the hashed input password and the hashed stored password
			if hashedPassword != storedPassword {
				errorMessage = "Invalid username or password."
			} else {
				// Compare if the accountID starts with "2", indicating an alumni account
				if accountID[:1] != "2" {
					errorMessage = "This is not an alumni account."
				} else {
					// Render next page if it's a valid alumni account
					http.Redirect(w, r, "/login/studentprofile", http.StatusSeeOther)
					return
				}
			}
		}

		// If there is an error, render the alumni page with the error message
		tmpl := template.Must(template.ParseFiles("../FrontEnd/LoginPage/alumnilog/alumni.html"))
		tmpl.Execute(w, struct {
			ErrorMessage string
			Username     string
			Password     string
		}{
			ErrorMessage: errorMessage,
			Username:     username,
			Password:     "", // Clear the password field
		})
		return
	}

	// ALUMNI page rendering
	tmpl := template.Must(template.ParseFiles("../FrontEnd/LoginPage/alumnilog/alumni.html"))
	tmpl.Execute(w, nil)
}
