package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"authentication/src/platform/auth"
	"authentication/src/platform/router"
)

func main() {
	godotenv.Load()

	miniAuthAddr := os.Getenv("MINIAUTH_ADDRESS")
	if miniAuthAddr == "" {
		miniAuthAddr = "localhost:50051"
	}

	authClient, err := auth.NewAuthClient(miniAuthAddr)
	if err != nil {
		log.Fatalf("Failed to create auth client: %v", err)
	}

	rtr := router.New(authClient)

	log.Print("HTTP server listening on http://localhost:3000/")
	if err := http.ListenAndServe("0.0.0.0:3000", rtr); err != nil {
		log.Fatalf("There was an error with the http server: %v", err)
	}
}
