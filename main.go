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

	// for n, repo := range loadGitRepositories(config) {
	// 	log.Printf("🔎 [%v] scanning %s", n, repo.path)
	// }

	startHttp(config)

	// print git directories
	// gitEntries, err := scanForGitDirs(config.root)
	// if err != nil {
	// 	log.Fatalf("Error scanning for GIT directories: %v", err)
	// }

	// save git directories to cvs file
	// saveToCsv(gitEntries, "repos.csv")
}
