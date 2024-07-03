package utils

import (
	"database/sql"
	"log"
)

func ConnectionToDB() *sql.DB {
	db, err := sql.Open("sqlite3", "../searchYardAPI/db/wan_show_main.db")

	if err != nil {
		log.Fatal(err)
	}

	return db
}
