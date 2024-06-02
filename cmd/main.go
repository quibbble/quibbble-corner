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
		log.Fatal("no 'PORT' env found")
	}
	adminUsername := os.Getenv("ADMIN_USERNAME")
	if adminUsername == "" {
		log.Fatal("no 'ADMIN_USERNAME' env found")
	}
	adminPassword := os.Getenv("ADMIN_PASSWORD")
	if adminPassword == "" {
		log.Fatal("no 'ADMIN_PASSWORD' env found")
	}

	qcorner.ServeHTTP(port, adminUsername, adminPassword)
}
