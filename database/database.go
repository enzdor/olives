package database

import (
	"database/sql"
	"fmt"
	"log"
)

func NewDB(user string, password string, name string) *sql.DB {
	cfg := fmt.Sprintf("%s:%s@tcp(127.0.0.1:3306)/%s?parseTime=true", user, password, name)

	db, err := sql.Open("mysql", cfg)
	if err != nil {
		log.Fatal(err)
	}

	return db
}
