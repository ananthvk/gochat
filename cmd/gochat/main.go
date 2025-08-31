package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ananthvk/gochat/internal"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/joho/godotenv"
)

func handlerHomePage(w http.ResponseWriter, r *http.Request) {
	log.Printf("GET /")
	fmt.Fprintf(w, "%s", "Hello there")
}

func main() {
	godotenv.Load()

	host, ok := os.LookupEnv("HOST")
	if !ok {
		host = "127.0.0.1"
	}
	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8000"
	}

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.HandleFunc("/", handlerHomePage)
	router.Mount("/api/v1/", internal.Routes())

	server := &http.Server{
		Addr:    host + ":" + port,
		Handler: router,
	}

	log.Printf("Server listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
