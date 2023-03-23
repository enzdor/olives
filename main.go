package main

import (
	"log"
	"net/http"
	"os"

	"github.com/jobutterfly/olives/controllers"
	"github.com/jobutterfly/olives/database"

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

	http.HandleFunc("/users/", h.GetOrDeleteUser)
	http.HandleFunc("/users", h.CreateUser)

	http.HandleFunc("/posts/", h.GetOrDeletePost)
	http.HandleFunc("/posts/subolive/", h.GetSubolivePosts)

	http.HandleFunc("/comments", h.CreateComment)

	log.Print("Listiening on port :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
