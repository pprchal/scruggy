package main

import (
	"log"
	"scruggy/config"
	"scruggy/git"
	"scruggy/http"
)

func main() {
	config.GlobalConfig = config.LoadConfiguration()

	// fetch status - todo: goroutine
	go func() {
		for n := range config.GlobalConfig.Repos {
			repo := config.GlobalConfig.Repos[n]
			status := git.Status(repo.Path)
			log.Printf("⌛ %d < git-status %s => %s", status.Status, repo.Path, status.Text)
		}
	}()

	// start gui
	http.StartHttp()
}
