package main

import (
	"encoding/json"
	"net/http"

	"github.com/56treskka/todo/internal/auth"
	"github.com/56treskka/todo/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUpdateItem(w http.ResponseWriter, r *http.Request) {
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

	itemIDString := r.PathValue("item_id")
	if itemIDString == "" {
		respondWithError(w, http.StatusBadRequest, "Missing item_id", nil)
		return
	}

	itemID, err := uuid.Parse(itemIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid item_id", err)
	}

	params := Paramaters{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding parameters", err)
		return
	}

	updatedItem, err := cfg.db.UpdateTask(r.Context(), database.UpdateTaskParams{
		Title:       params.Title,
		Description: params.Description,
		UserID:      userID,
		ID:          itemID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't update task", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Task{
		Id:          updatedItem.ID,
		Title:       updatedItem.Title,
		Description: updatedItem.Description,
	})

}
