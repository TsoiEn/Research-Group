package login

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"html/template"
	"net/http"
)

// Handler to fetch ADMIN account information
func adminHandler(w http.ResponseWriter, r *http.Request) {
	var errorMessage string

	if r.Method == http.MethodPost {
		// Parse form data
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Query the database for admin credentials
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
				// Compare if the accountID starts with "1", indicating an admin account
				if accountID[:1] != "1" {
					errorMessage = "This is not an admin account."
				} else {
					// Render next page if it's a valid admin account
					http.Redirect(w, r, "/login/adminstudentlist", http.StatusSeeOther)
					return
				}
			}
		}

		// If there is an error, render the admin page with the error message
		tmpl := template.Must(template.ParseFiles("../FrontEnd/LoginPage/adminlog/admin.html"))
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

	// ADMIN page rendering
	tmpl := template.Must(template.ParseFiles("../FrontEnd/LoginPage/adminlog/admin.html"))
	tmpl.Execute(w, nil)
}
