package main

import (
	"fmt"
	"log"
	"os"
)

func ScanStop() {

}

func newRepo(repo string) {
	// config.new_repos = nil

	for n := range config.entries {
		if config.entries[n].path == repo {
			log.Printf("🔎 repo already exists: %s", repo)
			return
		}
	}

	log.Printf("🔎 found repo: %s", repo)
	config.new_repos = append(config.new_repos, repo)
}

func ScanStart() {
	log.Printf("🔎 scanning %s", config.root)
	FindGitRepositories(config.root, newRepo)
}

func SyncAll() {

}

func AddRepo(repo string) {
	log.Printf("+ adding repo: %s", repo)

	f, err := os.OpenFile("config.ini", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = f.WriteString(fmt.Sprintf("\n[%s]", repo))
	if err != nil {
		log.Fatal(err)
	}
}
