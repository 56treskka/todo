package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/56treskka/todo/internal/database"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	db     *database.Queries
	secret string
}

func main() {
	dbURL := os.Getenv("DB_URL")
	secret := os.Getenv("SECRET")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
	dbQueries := database.New(db)

	apicfg := apiConfig{
		db:     dbQueries,
		secret: secret,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /register", apicfg.handlerRegister)
	mux.HandleFunc("POST /login", apicfg.handlerLogin)
	mux.HandleFunc("POST /create-item", apicfg.handlerCreateItem)
	mux.HandleFunc("POST /update-item/{item_id}", apicfg.handlerUpdateItem)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Fatal(server.ListenAndServe())
}
