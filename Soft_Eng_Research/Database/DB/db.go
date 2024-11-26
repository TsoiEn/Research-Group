package DB

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

// Initialize DB connection
func InitDB() (*sql.DB, error) {
	var err error
	db, err = sql.Open("mysql", "RaffyEL:WeMakingDBS1@E@tcp(127.0.0.1:3306)/stucredstorage")
	if err != nil {
		log.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
