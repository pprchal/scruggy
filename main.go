package main

import (
	"log"
)

func main() {
	// load config
	config := LoadConfiguration()

	//UpdateConfiguration(config)

	// repos := LoadGitConfig(config)
	for n := range config.entries {
		repo := config.entries[n]
		log.Printf("⌛ git-status %s", repo.path)
		// FetchGitStatus(repo)
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
