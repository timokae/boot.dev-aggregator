package main

import (
	"net/http"

	"github.com/timokae/boot.dev-aggregator/internal/auth"
	"github.com/timokae/boot.dev-aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (cfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Could not find api key")
			return
		}

		user, err := cfg.DB.FindUserByApiKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Could not find user")
			return
		}

		handler(w, r, user)
	}
}
