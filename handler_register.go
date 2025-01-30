package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/56treskka/todo/internal/auth"
	"github.com/56treskka/todo/internal/database"
)

func (cfg *apiConfig) handlerRegister(w http.ResponseWriter, r *http.Request) {
	type Paramaters struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	type Response struct {
		Token        string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}

	params := Paramaters{}
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Error decoding parameters", err)
		return
	}

	password, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't hash password", err)
		return
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		Name:     params.Name,
		Email:    params.Email,
		Password: password,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create user", err)
		return
	}

	token, err := auth.GenerateJWT(user.ID, cfg.secret, time.Duration(time.Minute*15))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate JWT", err)
		return
	}

	expirationTime := time.Hour * 24 * 7
	refreshTokenString, err := auth.GenerateRefreshToken(user.ID, expirationTime)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't generate refresh token", err)
	}

	refershToken, err := cfg.db.CreateRefreshToken(r.Context(), database.CreateRefreshTokenParams{
		UserID:    user.ID,
		Token:     refreshTokenString,
		ExpiresAt: time.Now().Add(expirationTime),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't create refresh token", err)
	}

	respondWithJSON(w, http.StatusCreated, Response{Token: token, RefreshToken: refershToken.Token})
}
