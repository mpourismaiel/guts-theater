package main

import (
	"log"

	"mpourismaiel.dev/guts/api"
)

func main() {
	log.Println("Starting project...")
	api.New("4000")
}
