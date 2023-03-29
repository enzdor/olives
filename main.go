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

	mux := http.NewServeMux()

	mux.HandleFunc("/users/", h.GetUser)
	mux.HandleFunc("/users/delete/", h.AuthMiddleware(h.DeleteUser))
	mux.HandleFunc("/users/create/", h.AuthMiddleware(h.CreateUser))

	mux.HandleFunc("/posts/create", h.AuthMiddleware(h.CreatePost))
	mux.HandleFunc("/posts/delete/", h.AuthMiddleware(h.AdminMiddleware(h.DeletePost)))
	mux.HandleFunc("/posts/", h.GetPost)
	mux.HandleFunc("/posts/subolive/", h.GetSubolivePosts)

	mux.HandleFunc("/comments/create", h.AuthMiddleware(h.CreateComment))
	mux.HandleFunc("/comments/delete", h.AuthMiddleware(h.AdminMiddleware(h.DeleteComment)))

	h.ExtenderMiddleware(mux)

	log.Print("Listiening on port :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
