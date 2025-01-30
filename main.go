package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"

	"github.com/56treskka/todo/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	db     *database.Queries
	secret string
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

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
	mux.HandleFunc("PUT /update-item/{item_id}", apicfg.handlerUpdateItem)
	mux.HandleFunc("GET /todos", apicfg.handlerGetItems)
	mux.HandleFunc("DELETE /delete-item/{item_id}", apicfg.handlerDeleteItem)
	mux.HandleFunc("POST /refresh", apicfg.handlerRefresh)
	mux.HandleFunc("POST /revoke", apicfg.handlerRevoke)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	log.Print("Serving files on port: 8080\n")
	log.Fatal(server.ListenAndServe())
}
