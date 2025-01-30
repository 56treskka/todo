package main

import (
	"net/http"
	"strconv"

	"github.com/56treskka/todo/internal/auth"
	"github.com/56treskka/todo/internal/database"
)

func (cfg *apiConfig) handlerGetItems(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Data  []Task `json:"data"`
		Page  int    `json:"page"`
		Limit int    `json:"limit"`
		Total int    `json:"total"`
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

	pageString := r.URL.Query().Get("page")
	if pageString == "" {
		pageString = "1"
	}
	page, err := strconv.Atoi(pageString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid page", err)
		return
	}

	limitString := r.URL.Query().Get("limit")
	if limitString == "" {
		limitString = "10"
	}
	limit, err := strconv.Atoi(limitString)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid limit", err)
		return
	}

	items, err := cfg.db.GetTasks(r.Context(), database.GetTasksParams{
		UserID: userID,
		Limit:  int32(limit),
		Offset: int32((page - 1) * limit),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get tasks", err)
		return
	}

	tasks := make([]Task, len(items))
	for i, item := range items {
		tasks[i] = Task{
			Id:          item.ID,
			Title:       item.Title,
			Description: item.Description,
		}
	}

	respondWithJSON(w, http.StatusOK, Response{
		Data:  tasks,
		Page:  page,
		Limit: limit,
		Total: len(tasks),
	})
}
