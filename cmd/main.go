package main

import (
	"log"
	"os"

	"github.com/quibbble/quibbble-controller/pkg/auth"
	"github.com/quibbble/quibbble-corner/internal/qcorner"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	authKey := os.Getenv("AUTH_KEY")

	// setup authenticate handler
	a, err := auth.NewAuth(authKey)
	if err != nil {
		log.Fatal(err)
	}

	qcorner.ServeHTTP(port, a.Retrieve)
}
