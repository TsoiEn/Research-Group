package login

import (
	"database/sql"
	"html/template"
	"net/http"

	sessionHandler "github.com/TsoiEn/Research-Group/Soft_Eng-Research/Soft_Eng_Research/Database/SessionStore"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Global variable for DB connection
var db *sql.DB

type Account struct {
	Username string
	Password string
}

type StudentInfo struct {
	StudentID string
	Course    string
	Fullname  string
	Email     string
}

type StudentList struct {
	StudentID string
	Course    string
	LastName  string
	FirstName string
	Email     string
}

// Handler to go into next page (TEST (will delete))
func successHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("success.html"))
	tmpl.Execute(w, nil)
}

// Handler to go into ADMIN STUDENT LIST PAGE and fetch all student's information
func adminStuListHandler(w http.ResponseWriter, r *http.Request) {
	// Query to fetch all student data
	query := `
		SELECT
			s.studentID AS studentID,
			s.course AS course,
			s.lname AS lastName,
			s.fname AS firstName,
			a.username AS email
		FROM
			students s
		INNER JOIN
			accounts a
		ON
			s.userID = a.accountID
		WHERE
			LEFT(a.accountID, 1) = '3'`

	// Execute query
	rows, err := db.Query(query)
	if err != nil {
		http.Error(w, "Failed to fetch student list", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Create a slice to hold the student data
	var students []StudentList

	// Loop through rows and populate the slice
	for rows.Next() {
		var student StudentList
		if err := rows.Scan(&student.StudentID, &student.Course, &student.LastName, &student.FirstName, &student.Email); err != nil {
			http.Error(w, "Failed to scan student data", http.StatusInternalServerError)
			return
		}

		// Append the student to the list
		students = append(students, student)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		http.Error(w, "Error while fetching data", http.StatusInternalServerError)
		return
	}

	// Render the AdminStudentList with the fetched data
	tmpl := template.Must(template.ParseFiles("../FrontEnd/HomePage/AdminHomePage/AdminStudentList/AdminStudentList.html"))
	tmpl.Execute(w, students)
}

// Handler to go into STUDENT PROFILE PAGE and fetch the specific student's information
func stuProfHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve session
	session, _ := sessionHandler.StoreStuProf.Get(r, "student-session")
	accountID, ok := session.Values["accountID"].(string)
	if !ok || accountID == "" {
		http.Error(w, "Unauthorized: Account ID is missing", http.StatusUnauthorized)
		return
	}

	// Fetch student profile data using the accountID
	student, err := fetchStudentProfile(accountID)
	if err != nil {
		http.Error(w, "Failed to fetch profile data", http.StatusInternalServerError)
		return
	}

	// Render the StudentProfile.html with the fetched data
	tmpl := template.Must(template.ParseFiles("../FrontEnd/HomePage/StudentHomePage/StudentProfile/StudentProfile.html"))
	tmpl.Execute(w, student)
}

// Handler to fetch specific student information from database
func fetchStudentProfile(accountID string) (StudentInfo, error) {
	var student StudentInfo

	// Query to fetch details by joining account and students tables
	query := `
		SELECT
			s.studentID,
			s.course,
			CONCAT(s.fname, ' ', s.lname) AS fullname,
			a.username
		FROM
			students s
		INNER JOIN
			accounts a
		ON
			s.userID = a.accountID
		WHERE
			a.accountID = ?`

	// Execute query and scan results into the struct
	err := db.QueryRow(query, accountID).Scan(
		&student.StudentID,
		&student.Course,
		&student.Fullname,
		&student.Email,
	)

	return student, err
}

func MainLogin(database *sql.DB) *mux.Router {
	db = database // Assign the DB connection to the global variable

	// Set up routes (r)
	r := mux.NewRouter()

	// ROUTES
	r.HandleFunc("/adminlogin", adminHandler).Methods("GET", "POST")      // Admin login page
	r.HandleFunc("/alumnilogin", alumniHandler).Methods("GET", "POST")    // Alumni login page
	r.HandleFunc("/studentlogin", studentHandler).Methods("GET", "POST")  // Student login page
	r.HandleFunc("/success", successHandler).Methods("GET")               // TESTING
	r.HandleFunc("/adminstudentlist", adminStuListHandler).Methods("GET") // To ADMIN STUDENT LIST page
	r.HandleFunc("/studentprofile", stuProfHandler).Methods("GET")        // To STUDENT PROFILE page

	// STATIC FILES (for JS, images, etc.)
	fsAdminStudentList := http.FileServer(http.Dir("../FrontEnd/HomePage/AdminHomePage/AdminStudentList"))
	r.PathPrefix("/AdminStudentList/").Handler(http.StripPrefix("/AdminStudentList", fsAdminStudentList))

	fsStudentProfile := http.FileServer(http.Dir("../FrontEnd/HomePage/StudentHomePage/StudentProfile"))
	r.PathPrefix("/StudentProfile/").Handler(http.StripPrefix("/StudentProfile", fsStudentProfile))

	fsAdmin := http.FileServer(http.Dir("../FrontEnd/LoginPage/adminlog"))
	r.PathPrefix("/adminlog/").Handler(http.StripPrefix("/adminlog", fsAdmin))

	fsAlumni := http.FileServer(http.Dir("../FrontEnd/LoginPage/alumnilog"))
	r.PathPrefix("/alumnilog/").Handler(http.StripPrefix("/alumnilog", fsAlumni))

	fsStud := http.FileServer(http.Dir("../FrontEnd/LoginPage"))
	r.PathPrefix("/").Handler(http.StripPrefix("/", fsStud))

	return r
}
