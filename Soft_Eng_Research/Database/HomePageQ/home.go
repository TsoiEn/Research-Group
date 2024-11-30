package home

import (
	"database/sql"
	"encoding/base64"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"

	sessionHandler "Soft_Eng_Research/Database/SessionStore"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Global variable for DB connection
var db *sql.DB

type Credential struct {
	ID         string
	FileData   string
	FileType   string
	Type       string
	Issuer     string
	DateIssued string
	Status     string
}

// FOR ADMIN
// Handler to go into ADMIN CREDENTIALS PAGE
func adminCredHandler(w http.ResponseWriter, r *http.Request) {
	// Get the studentID from query parameters
	studentID := r.URL.Query().Get("studentID")
	if studentID == "" {
		http.Error(w, "Student ID is required", http.StatusBadRequest)
		return
	}

	log.Print("Opening account of ", "[", studentID, "]")

	// Query the database for the student's credentials
	credentials, err := fetchStudentCredentialsByStudentID(studentID)
	if err != nil {
		http.Error(w, "Failed to fetch credentials", http.StatusInternalServerError)
		return
	}

	// Render next page
	tmpl := template.Must(template.ParseFiles("../FrontEnd/HomePage/AdminHomePage/AdminCredentials/AdminCredentials.html"))
	tmpl.Execute(w, credentials)
}

func fetchStudentCredentialsByStudentID(studentID string) (map[string][]map[string]string, error) {
	query := `
		SELECT c.credentialID, c.filedata, c.filetype, c.type, c.issuer, c.date_issued, c.status
		FROM credentials c
		INNER JOIN accounts a ON c.ownerID = a.accountID
		INNER JOIN students s ON a.accountID = s.userID
		WHERE s.studentID = ?`

	rows, err := db.Query(query, studentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	credentials := make(map[string][]map[string]string)

	for rows.Next() {
		var credentialID, filetype, credentialType, issuer, dateIssued, status string
		var filedata []byte
		if err := rows.Scan(&credentialID, &filedata, &filetype, &credentialType, &issuer, &dateIssued, &status); err != nil {
			return nil, err
		}

		credential := map[string]string{
			"credentialID": credentialID,
			"filedata":     base64.StdEncoding.EncodeToString(filedata),
			"filetype":     filetype,
			"issuer":       issuer,
			"dateIssued":   dateIssued,
			"status":       status,
		}

		credentials[credentialType] = append(credentials[credentialType], credential)
	}

	return credentials, nil
}

func addCredential(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form to retrieve file data
	err := r.ParseMultipartForm(10 << 20) // Limit file size to 10 MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("filedata")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	// Get file type
	filetype := handler.Header.Get("Content-Type")

	// Retrieve studentID from the request form
	studentID := r.FormValue("studentID")
	if studentID == "" {
		http.Error(w, "Student ID is required", http.StatusBadRequest)
		return
	}

	// Fetch ownerID from the accounts table
	var ownerID string
	query := `
		SELECT a.accountID
		FROM accounts a
		INNER JOIN students s ON s.userID = a.accountID
		WHERE s.studentID = ?`

	err = db.QueryRow(query, studentID).Scan(&ownerID)
	if err != nil {
		http.Error(w, "Error fetching account ID", http.StatusInternalServerError)
		return
	}

	var issuer string = "Admin"

	// Get the credential type from the form
	credentialType := r.FormValue("type")
	if credentialType != "academic" && credentialType != "non-academic" && credentialType != "certificate" {
		http.Error(w, "Invalid credential type", http.StatusBadRequest)
		return
	}

	// Set date format
	dateIssued := time.Now().Format("2006-01-02")

	// Insert data into the database
	query = `
	INSERT INTO credentials (ownerID, filedata, filetype, type, issuer, date_issued, status)
	VALUES (?, ?, ?, ?, ?, ?, 'active')`

	log.Printf("ownerID: %s, fileBytes: %d bytes, filetype: %s, credentialType: %s, issuer: %s, dateIssued: %s",
		ownerID, len(fileBytes), filetype, credentialType, issuer, dateIssued)

	_, err = db.Exec(query, ownerID, fileBytes, filetype, credentialType, issuer, dateIssued)
	if err != nil {
		log.Printf("Error executing query: %v", err)
		http.Error(w, "Error saving to database", http.StatusInternalServerError)
		return
	}

	// Redirect back to the admin student list page
	http.Redirect(w, r, "/login/adminstudentlist", http.StatusSeeOther)
}

func addCredentialPageHandler(w http.ResponseWriter, r *http.Request) {
	// Assuming you have a session or query parameter for the current studentID
	studentID := r.URL.Query().Get("studentID")
	if studentID == "" {
		http.Error(w, "Student ID is required", http.StatusBadRequest)
		return
	}

	// Render the page with studentID passed to the template
	tmpl := template.Must(template.ParseFiles("../FrontEnd/HomePage/AdminHomePage/AdminCredentials/AdminCredentials.html"))
	tmpl.Execute(w, map[string]interface{}{
		"StudentID": studentID, // Pass studentID to template
	})
}

// FOR STUDENT
// Handler to go into STUDENT CREDENTIALS PAGE
func stuCredHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := sessionHandler.StoreStuProf.Get(r, "student-session")
	accountID, ok := session.Values["accountID"].(string)
	if !ok || accountID == "" {
		http.Error(w, "Unauthorized: Account ID is missing", http.StatusUnauthorized)
		return
	}

	credentials, err := fetchStudentCredentials(accountID)
	if err != nil {
		http.Error(w, "Failed to fetch credentials", http.StatusInternalServerError)
		return
	}

	// Render next page
	tmpl := template.Must(template.ParseFiles("../FrontEnd/HomePage/StudentHomePage/StudentCredentials/StudentCred.html"))
	tmpl.Execute(w, credentials)
}

// Handler to fetch credentials for the logged-in student based on their accountID
// Function to fetch credentials for a student
func fetchStudentCredentials(accountID string) (map[string][]map[string]string, error) {
	query := `
        SELECT credentialID, filedata, filetype, type, issuer, date_issued, status
        FROM credentials
        WHERE ownerID = ?`

	rows, err := db.Query(query, accountID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	credentials := make(map[string][]map[string]string)

	for rows.Next() {
		var credentialID, filetype, credentialType, issuer, dateIssued, status string
		var filedata []byte
		if err := rows.Scan(&credentialID, &filedata, &filetype, &credentialType, &issuer, &dateIssued, &status); err != nil {
			return nil, err
		}

		credential := map[string]string{
			"credentialID": credentialID,
			"filedata":     base64.StdEncoding.EncodeToString(filedata),
			"filetype":     filetype,
			"issuer":       issuer,
			"dateIssued":   dateIssued,
			"status":       status,
		}

		credentials[credentialType] = append(credentials[credentialType], credential)
	}

	return credentials, nil
}

func addNonAcademicHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Parse the form to retrieve file data
	err := r.ParseMultipartForm(10 << 20) // Limit file size to 10 MB
	if err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}

	file, handler, err := r.FormFile("filedata")
	if err != nil {
		http.Error(w, "Error retrieving file", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	// Read file content
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	// Get file type
	filetype := handler.Header.Get("Content-Type")

	// Get session
	session, _ := sessionHandler.StoreStuProf.Get(r, "student-session")

	// Get ownerID
	ownerID, ok := session.Values["accountID"].(string)
	if !ok || ownerID == "" {
		http.Error(w, "Unauthorized: Unable to get user ID", http.StatusUnauthorized)
		return
	}

	// Get name the of issuer (the student)
	issuer, ok := session.Values["studentName"].(string)
	if !ok || issuer == "" {
		// Fallback: Fetch from database
		err := db.QueryRow("SELECT CONCAT(fname, ' ', lname) FROM students WHERE userID = ?", ownerID).Scan(&issuer)
		if err != nil {
			http.Error(w, "Error fetching student name", http.StatusInternalServerError)
			return
		}
	}

	// Set data format
	dateIssued := time.Now().Format("2006-01-02")

	// Insert data into the database
	query := `
        INSERT INTO credentials (ownerID, filedata, filetype, type, issuer, date_issued, status)
        VALUES (?, ?, ?, 'non-academic', ?, ?, 'active')`

	_, err = db.Exec(query, ownerID, fileBytes, filetype, issuer, dateIssued)
	if err != nil {
		http.Error(w, "Error saving to database", http.StatusInternalServerError)
		return
	}

	// Redirect back to the credentials page
	http.Redirect(w, r, "/home/studentcredentials", http.StatusSeeOther)
}

func MainHome(database *sql.DB) *mux.Router {
	db = database // Assign the DB connection to the global variable

	// Set up routes (r)
	r := mux.NewRouter()

	// ROUTES
	r.HandleFunc("/admincredentials", adminCredHandler).Methods("GET") // To ADMIN CREDENTIALS PAGE
	r.HandleFunc("/studentcredentials", stuCredHandler).Methods("GET") // To STUDENT CREDENTIALS PAGE

	// ENDPOINT ROUTES
	r.HandleFunc("/add-non-academic", addNonAcademicHandler)
	r.HandleFunc("/add-credential", addCredential)
	r.HandleFunc("/add-credential", addCredentialPageHandler)

	// r.HandleFunc("/admin-add-academic", a)

	// STATIC FILES (Serve CSS, JS, images)
	fsAdminCred := http.FileServer(http.Dir("../FrontEnd/HomePage/AdminHomePage/AdminCredentials"))
	r.PathPrefix("/AdminCredentials/").Handler(http.StripPrefix("/AdminCredentials", fsAdminCred))

	fsStuCred := http.FileServer(http.Dir("../FrontEnd/HomePage/StudentHomePage/StudentCredentials"))
	r.PathPrefix("/StudentCredentials/").Handler(http.StripPrefix("/StudentCredentials", fsStuCred))

	return r
}
