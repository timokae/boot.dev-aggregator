package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/timokae/boot.dev-aggregator/internal/database"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("CONN")
	if dbURL == "" {
		log.Fatalln("Database connection string missing")
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatalln(err)
	}

	dbQueries := database.New(db)
	cfg := apiConfig{
		DB: dbQueries,
	}

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    "localhost:" + port,
		Handler: mux,
	}

	mux.HandleFunc("GET /v1/healthz", handlerReadiness)
	mux.HandleFunc("GET /v1/err", handlerErr)

	mux.HandleFunc("POST /v1/users", cfg.handlerUsersCreate)
	mux.HandleFunc("GET /v1/users", cfg.middlewareAuth(cfg.handlerUserGet))

	mux.HandleFunc("POST /v1/feeds", cfg.middlewareAuth(cfg.handlerFeedsCreate))
	mux.HandleFunc("GET /v1/feeds", cfg.handlerFeedsGet)

	mux.HandleFunc("POST /v1/feed_follows", cfg.middlewareAuth(cfg.handlerFeedFollowsCreate))
	mux.HandleFunc("DELETE /v1/feed_follows/{id}", cfg.middlewareAuth(cfg.handlerFeedFollowsDelete))
	mux.HandleFunc("GET /v1/feed_follows", cfg.middlewareAuth(cfg.handlerFeedFollowsGet))

	mux.HandleFunc("GET /v1/posts", cfg.middlewareAuth(cfg.handlerGetPostsForUser))

	go scrapeFeeds(cfg.DB, 10, time.Minute)

	log.Printf("Serving on port: %s\n", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
