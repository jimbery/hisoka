package main

import (
	"fmt"
	"hisoka/internal/httpserver"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Errorf("failed to load environment variables: %s", err)
	}

	httpserver.Listen()
}
