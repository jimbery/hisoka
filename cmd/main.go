package main

import (
	"flag"
	"hisoka/internal/httpserver"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	// Check if the environment is "local" or "development"
	localFlag := flag.Bool("local", false, "Set to load .env for local development")

	flag.Parse()

	// If the --local flag is set, load the .env file
	if *localFlag {
		err := godotenv.Load()
		if err != nil {
			log.Fatal("Error loading .env file")
			return
		}
	}

	httpserver.Listen()
}
