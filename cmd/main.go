package main

import (
	"fmt"
	"hisoka/internal/httpclient"
	"log"
)

func main() {
	var anime httpclient.AnimeDetails
	taskCompleted := make(chan bool)

	searchTerm := "hum"

	// Start the API req in a goroutine
	go func() {
		var err error
		anime, err = httpclient.SearchAnime(searchTerm, &httpclient.RealHTTPClient{})
		if err != nil {
			log.Fatalln(err)
		}

		// Signal that the task has completed by sending true to the channel
		taskCompleted <- true
	}()

	// Block until the task is completed
	<-taskCompleted

	fmt.Println(anime)
}
