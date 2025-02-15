package main

import (
	"fmt"
	"log"
	"os"
)

func RepoAction(repo string, action string, remote string) {
	gitRepo := findRepository(repo)
	if gitRepo == nil {
		log.Fatalf("😭 repo not found: %s", repo)
		return
	}

	fmt.Printf("git[%s] %s %s", remote, action, repo)
}

func ScanStop() {
}

func ScanNewRepo(repo string) {
	gitRepo := findRepository(repo)
	if gitRepo == nil {
		log.Printf("🔎 found repo: %s", repo)
		config.new_repos = append(config.new_repos, repo)
		return
	}

	log.Printf("🔎 repo already exists: %s", repo)
}

func findRepository(repo string) *GitRepo {
	for n := range config.repos {
		if config.repos[n].path == repo {
			return &config.repos[n]
		}
	}
	return nil
}

func ScanStart() {
	log.Printf("🔎 scanning %s", config.root)
	FindGitRepositories(config.root, ScanNewRepo)
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

	actions := "push-a,pull-a"
	_, err = f.WriteString(fmt.Sprintf("\n[%s]\nactions=%s", repo, actions))
	if err != nil {
		log.Fatal(err)
	}
}
