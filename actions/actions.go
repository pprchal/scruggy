package actions

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"scruggy/config"
	"scruggy/git"
)

func OpenTerminalWindow(path string) {
	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", path).Start()
	case "windows":
		// TODO: maybe start is better...
		exec.Command("rundll32", "url.dll,FileProtocolHandler", path).Start()
	case "darwin":
		exec.Command("open", path).Start()
	default:
		log.Fatalf("😭 unsupported platform: %s", runtime.GOOS)
	}
}

func Quit() {
	log.Printf("bye!")
	os.Exit(0)
}

func repoAction(repo string, action string, remote string) {
	gitRepo := findRepository(repo)
	if gitRepo == nil {
		log.Fatalf("😭 repo not found: %s", repo)
		return
	}

	var result git.GitResult
	switch action {
	case "push":
		result = git.Push(repo, remote)

	case "pull":
		result = git.Pull(repo, remote)
	}

	log.Printf("%d < git %s %s [%s] => %s", result.Status, action, remote, repo, result.Text)
}

func RepoActions(repo string, gitActions string) {
	for _, action := range config.ParseActions(gitActions) {
		repoAction(repo, action.Action, action.Remote)
	}
}

func Refresh() {
	config.GlobalConfig = config.LoadConfiguration()
	for n := range config.GlobalConfig.Repos {
		repo := &config.GlobalConfig.Repos[n]
		status := git.Status(repo.Path)

		if status.Text == "" {
			repo.Status = 0
		} else {
			repo.Status = 1
		}
		log.Printf("⌛ %d < git-status %s => %s", status.Status, repo.Path, status.Text)
	}
}

func ScanNewRepo(path string) {
	gitRepo := findRepository(path)
	if gitRepo == nil {
		log.Printf("🔎 found: %s", path)
		config.GlobalConfig.NewRepos = append(config.GlobalConfig.NewRepos, path)
		return
	}

	log.Printf("🔎 already exists: %s", path)
}

func findRepository(repo string) *git.GitRepo {
	for n := range config.GlobalConfig.Repos {
		if config.GlobalConfig.Repos[n].Path == repo {
			return &config.GlobalConfig.Repos[n]
		}
	}
	return nil
}

func ScanStart() {
	log.Printf("🔎 scanning %s", config.GlobalConfig.Root)
	git.FindGitRepositories(config.GlobalConfig.Root, ScanNewRepo)
}

func SyncAll() {
	for _, repo := range config.GlobalConfig.Repos {
		for _, action := range repo.Actions {
			repoAction(repo.Path, action.Action, action.Remote)
		}
	}
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
	repo := git.GitRepo{Path: path}
	git.LoadGitConfig(&repo)

	// default actions
	actions := ""
	for _, remote := range repo.Remotes {
		actions += fmt.Sprintf("push-%s,pull-%s", remote.Name, remote.Name)
	}

	_, err = f.WriteString(fmt.Sprintf("\n[%s]\nactions=%s\n", repo.Path, actions))
	if err != nil {
		log.Fatal(err)
	}

	config.GlobalConfig.Repos = append(config.GlobalConfig.Repos, repo)
}
