package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/56treskka/todo/internal/auth"
)

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type Paramaters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type Response struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	decoder := json.NewDecoder(r.Body)
	params := Paramaters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding parameters", err)
		return
	}

	user, err := cfg.db.GetUser(r.Context(), params.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user", err)
		return
	}
	if err := auth.ComparePassword(user.Password, params.Password); err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid password", err)
		return
	}

	accessToken, err := auth.GenerateJWT(user.ID, cfg.secret, time.Minute*15)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate access token", err)
		return
	}
	refreshToken, err := auth.GenerateRefreshToken(user.ID, time.Hour*24*7)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate refresh token", err)
		return
	}

	respondWithJSON(w, http.StatusCreated, Response{
		Token:        accessToken,
		RefreshToken: refreshToken,
	})
}
