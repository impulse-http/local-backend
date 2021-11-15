package main

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/impulse-http/local-backend/pkg/database"
	"github.com/impulse-http/local-backend/pkg/service"
	"github.com/impulse-http/local-backend/pkg/service/requests"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {
	db, err := sql.Open("sqlite3", "dist/impulse.db")
	if err != nil {
		log.Fatal("cannot create sqlite3 db")
	}
	defer db.Close()
	dbWrapper := database.NewDatabase(db)

	s := service.NewService(dbWrapper)
	r := mux.NewRouter()

	requests.AddRequestHandlers(s, r)

	log.Println("Running at localhost:8090")
	http.Handle("/", r)
	_ = http.ListenAndServe(":8090", nil)
}
