package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}

	port := os.Getenv("PORT")
	if err != nil {
		port = "8080"
	}

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    "localhost:" + port,
		Handler: mux,
	}

	mux.HandleFunc("GET /v1/healthz", func(w http.ResponseWriter, r *http.Request) {
		respondWithJSON(w, http.StatusOK, struct {
			Status string
		}{
			Status: "ok",
		})
	})

	mux.HandleFunc("GET /v1/err", func(w http.ResponseWriter, r *http.Request) {
		respondWithError(w, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	})

	log.Printf("Serving on port: %s\n", port)
	err = server.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
