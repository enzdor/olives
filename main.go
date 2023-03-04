package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jobutterfly/olives/database"
	"github.com/jobutterfly/olives/controllers"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

func main() {
	errEnv := godotenv.Load()
	if errEnv != nil {
	    log.Fatal(errEnv)
	}
	user := os.Getenv("DBUSER")
	pass := os.Getenv("DBPASS")
	name := os.Getenv("DBNAME")

	db := database.NewDB(user, pass, name)
	h := controllers.NewHandler(db, "jwtkey")

	http.HandleFunc("/users/", h.GetUser)

	log.Print("Listiening on port :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
