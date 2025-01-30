package main

import (
	"encoding/json"
	"net/http"

	"github.com/56treskka/todo/internal/auth"
	"github.com/56treskka/todo/internal/database"
)

func (cfg *apiConfig) handlerCreateItem(w http.ResponseWriter, r *http.Request) {
	type Paramaters struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}

	tokenString, err := auth.GetBearerToken(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	userID, err := auth.ValidateJWT(tokenString, cfg.secret)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	params := Paramaters{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding parameters", err)
		return
	}

	task, err := cfg.db.CreateTask(r.Context(), database.CreateTaskParams{
		Title:       params.Title,
		Description: params.Description,
		UserID:      userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create task", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Task{
		Id:          task.ID,
		Title:       task.Title,
		Description: task.Description,
	})
}
