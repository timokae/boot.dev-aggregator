package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/timokae/boot.dev-aggregator/internal/database"
)

func (cfg *apiConfig) handlerFeedFollowsCreate(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID string `json:"feed_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not decode parameters")
		return
	}

	feedID, err := uuid.Parse(params.FeedID)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusBadRequest, "Invalid feed id")
		return
	}

	feed, err := cfg.DB.GetFeed(r.Context(), feedID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "Could not find feed")
		return
	}

	feedFollow, err := cfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		UserID:    user.ID,
		FeedID:    feed.ID,
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
	})
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Could not create feed follow")
		return
	}

	respondWithJSON(w, http.StatusCreated, databaseFeedFollowToFeedFollow(feedFollow))
}

func (cfg *apiConfig) handlerFeedFollowsDelete(w http.ResponseWriter, r *http.Request, user database.User) {
	idStr := r.PathValue("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not parse ID")
		return
	}

	err = cfg.DB.DeleteFeedFollow(r.Context(), id)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusNotFound, "Could not find or delete feed follow")
		return
	}

	respondWithJSON(w, http.StatusOK, struct{}{})
}

func (cfg *apiConfig) handlerFeedFollowsGet(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollows, err := cfg.DB.GetFeedFollowsOfUser(r.Context(), user.ID)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, "Could not get feed follows")
		return
	}

	respondWithJSON(w, http.StatusOK, databaseFeedFollowsToFeedFollows(feedFollows))
}
