package main

import (
	"log"
)

var config Configuration

func main() {
	config = LoadConfiguration()

	// fetch status - todo: goroutine
	for n := range config.repos {
		repo := config.repos[n]
		log.Printf("⌛ git-status %s", repo.path)
	}

	// start gui
	startHttp(config)
}
