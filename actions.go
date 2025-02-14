package main

import "log"

func ScanStop() {

}

func newRepo(repo string) {
	log.Printf("+ found repo: %s", repo)
	// config.new_repos = nil
	config.new_repos = append(config.new_repos, repo)
}

func ScanStart() {
	log.Printf("🔎 scanning %s", config.root)
	FindGitRepositories(config.root, newRepo)
}

func SyncAll() {

}
