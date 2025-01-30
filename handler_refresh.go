package main

import (
	"net/http"
	"time"

	"github.com/56treskka/todo/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Token string `json:"token"`
	}
	refreshToken, err := auth.GetBearerToken(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	userID, err := cfg.db.GetUserFromRefreshToken(r.Context(), refreshToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	tokenString, err := auth.GenerateJWT(userID, cfg.secret, time.Minute*15)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate JWT", err)
		return
	}

	respondWithJSON(w, http.StatusOK, Response{
		Token: tokenString,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Unauthorized", err)
		return
	}

	if err = cfg.db.RevokeRefreshToken(r.Context(), refreshToken); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't revoke refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusNoContent, nil)
}
