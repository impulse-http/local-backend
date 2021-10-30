package main

import (
	"database/sql"
	"github.com/impulse-http/local-backend/pkg/service"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

func main() {
	db, err := sql.Open("sqlite3", "./impulse.db")
	if err != nil {
		log.Fatal("cannot create sqlite3 db")
	}
	defer db.Close()
	s := service.NewService(db)

	http.HandleFunc("/makeRequest", s.MakeRequestHandler)

	http.ListenAndServe(":8090", nil)
}
