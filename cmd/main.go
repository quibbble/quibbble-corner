package main

import (
	"log"
	"os"

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

	qcorner.ServeHTTP(port)
}
