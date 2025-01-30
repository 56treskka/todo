package main

import (
	"net/http"

	"github.com/56treskka/todo/internal/auth"
	"github.com/56treskka/todo/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerDeleteItem(w http.ResponseWriter, r *http.Request) {
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

	tokenIDString := r.PathValue("item_id")
	if tokenIDString == "" {
		respondWithError(w, http.StatusBadRequest, "Missing item_id", nil)
		return
	}

	tokenID, err := uuid.Parse(tokenIDString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid item_id", err)
		return
	}

	if err = cfg.db.DeleteTask(r.Context(), database.DeleteTaskParams{
		ID:     tokenID,
		UserID: userID,
	}); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't delete task", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
