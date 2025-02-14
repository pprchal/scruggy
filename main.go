package main

import (
	"log"
)

var config Configuration

func main() {
	// load config
	config = LoadConfiguration()
	// config.taskMessages = make(chan string)

	for n := range config.entries {
		repo := config.entries[n]
		log.Printf("⌛ git-status %s", repo.path)
	}

	startHttp(config)
}
