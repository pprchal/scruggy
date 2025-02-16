package main

import (
	"fmt"
	"log"
	"os"
)

func Quit() {
	log.Printf("bye!")
}

func RepoAction(repo string, action string, remote string) {
	gitRepo := findRepository(repo)
	if gitRepo == nil {
		log.Fatalf("😭 repo not found: %s", repo)
		return
	}

	var result GitResult
	switch action {
	case "push":
		result = gitPush(repo, remote)

	case "pull":
		result = gitPull(repo, remote)
	}

	log.Printf("%d < git %s %s [%s] => %s", result.status, action, remote, repo, result.text)
}

func ScanStop() {
}

func ScanNewRepo(path string) {
	gitRepo := findRepository(path)
	if gitRepo == nil {
		log.Printf("🔎 found: %s", path)
		config.new_repos = append(config.new_repos, path)
		return
	}

	log.Printf("🔎 already exists: %s", path)
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

func AddRepo(path string) {
	log.Printf("+ adding repo: %s", path)

	// open ini file
	f, err := os.OpenFile("config.ini", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	// load git config
	repo := GitRepo{path: path}
	LoadGitConfig(&repo)

	// default actions
	actions := ""
	for _, remote := range repo.remotes {
		actions += fmt.Sprintf("push-%s,pull-%s", remote.name, remote.name)
	}

	_, err = f.WriteString(fmt.Sprintf("\n[%s]\nactions=%s\n", repo.path, actions))
	if err != nil {
		log.Fatal(err)
	}
}
