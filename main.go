package main

import (
	"log"
	"net/http"
	"poll-app/handler"
	"poll-app/storage/postgres"
	"time"

	"github.com/gorilla/schema"
	"github.com/gorilla/sessions"
)

func main() {

	session := sessions.NewCookieStore([]byte("my_secret"))

	newDbString := newDBFromConfig()

	store, err := postgres.NewStorage(newDbString)

	if err != nil {
		log.Println("error db")
	}
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	r, err := handler.NewServer(store, decoder, session)

	if err != nil {
		log.Println("error on handelr")
	}

	srv := &http.Server{
		Handler: r,
		Addr:    "127.0.0.1:8000",
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())

}

func newDBFromConfig() string {
	dbParams := " " + "user=polluser"
	dbParams += " " + "host=localhost"
	dbParams += " " + "port=5432"
	dbParams += " " + "dbname=polldb"
	dbParams += " " + "password=admin"
	dbParams += " " + "sslmode=disable"

	return dbParams
}
