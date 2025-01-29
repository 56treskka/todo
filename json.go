package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func respondWithError(w http.ResponseWriter, code int, msg string, err error) {
	if err != nil {
		log.Println(err)
	}
	type ErrorResponse struct {
		Error string `json:"error"`
	}
	respondWithJSON(w, code, ErrorResponse{
		Error: msg,
	})
}

func respondWithJSON(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	data, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
		w.WriteHeader(500)
		return
	}
	w.WriteHeader(code)
	w.Write(data)
}
