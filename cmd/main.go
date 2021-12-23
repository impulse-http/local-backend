package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/impulse-http/local-backend/pkg/database"
	"github.com/impulse-http/local-backend/pkg/service"
	"github.com/impulse-http/local-backend/pkg/service/collections"
	"github.com/impulse-http/local-backend/pkg/service/requests"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
	"log"
	"net/http"
)

func main() {
	db, err := sql.Open("sqlite3", "dist/impulse.db")
	if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		log.Fatal(err)
	}
	if err != nil {
		log.Fatal("cannot create sqlite3 db")
	}
	defer db.Close()
	dbWrapper := database.NewDatabase(db)

	s := service.NewService(dbWrapper)
	r := mux.NewRouter()
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowCredentials: true,
	})

	requests.AddRequestHandlers(s, r)
	collections.AddCollectionsHandlers(s, r)

	log.Println("Running at localhost:8090")
	http.Handle("/", r)
	handler := c.Handler(r)
	_ = http.ListenAndServe(":8090", handler)
}
