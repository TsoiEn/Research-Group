package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	blockchain "github.com/TsoiEn/Research-Group/Soft_Eng_Research/Blockchain_Core"
	dbHandler "github.com/TsoiEn/Research-Group/Soft_Eng_Research/Database/DB"
	homeHandler "github.com/TsoiEn/Research-Group/Soft_Eng_Research/Database/HomePageQ"
	loginHandler "github.com/TsoiEn/Research-Group/Soft_Eng_Research/Database/LoginPageQ"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Global variable for DB connection
var db *sql.DB

var blockchainCore *blockchain.Blockchain

func initBlockchain() *blockchain.Blockchain {
	nodeID := "node1"                   // Example node ID
	peers := []string{"node2", "node3"} // Example peer nodes

	blockchainCore := blockchain.NewBlockchain(nodeID, peers)
	return blockchainCore
}

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

	log.Println("Initializing blockchain...")

	bc := blockchain.NewBlockchain("node1", []string{"node2", "node3"})
	log.Printf("Blockchain initialized: %+v", bc)

}
