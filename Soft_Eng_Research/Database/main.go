package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	dbHandler "github.com/TsoiEn/Research-Group/Soft_Eng-Research/Soft_Eng_Research/Database/DB"
	homeHandler "github.com/TsoiEn/Research-Group/Soft_Eng-Research/Soft_Eng_Research/Database/HomePageQ"
	loginHandler "github.com/TsoiEn/Research-Group/Soft_Eng-Research/Soft_Eng_Research/Database/LoginPageQ"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Global variable for DB connection
var db *sql.DB

func main() {
	// Initialize the database
	db, err := dbHandler.InitDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Initialize routers from different packages
	rLogin := loginHandler.MainLogin(db)
	rHome := homeHandler.MainHome(db)

	// Combine the routers using a parent router
	r := mux.NewRouter()

	// Mount routes from different packages to different URL prefixes
	r.PathPrefix("/login").Handler(http.StripPrefix("/login", rLogin)) // Mount login routes under "/login"
	r.PathPrefix("/home").Handler(http.StripPrefix("/home", rHome))    // Mount home routes under "/home"

	// Start the local server
	fmt.Println("Server is running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))

}
