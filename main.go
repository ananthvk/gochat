package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

func homePageHandler(w http.ResponseWriter, r *http.Request) {
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

	mux := http.NewServeMux()
	mux.HandleFunc("/", homePageHandler)

	server := &http.Server{
		Addr:    host + ":" + port,
		Handler: mux,
	}

	log.Printf("Server listening on %s", server.Addr)
	log.Fatal(server.ListenAndServe())
}
