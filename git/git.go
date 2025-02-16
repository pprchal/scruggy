package git

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
)

// git remote repo - .git/config
type GitRemote struct {
	Name string
	Url  string
}

// config-action pair - .config.ini
type GitAction struct {
	Remote string
	Action string
}

// git repository - .config.ini
type GitRepo struct {
	Path    string
	Remotes []GitRemote
	Actions []GitAction
	State   string
}

type GitResult struct {
	Text   string
	Status int
}

func FindGitRepositories(root string, f func(string)) {
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() && info.Name() == ".git" {
			path = strings.ReplaceAll(path, ".git", "")
			path = strings.TrimSuffix(path, string(filepath.Separator))
			f(path)
		}

		return nil
	})

	if err != nil {
		log.Fatalf("😭 error scanning for GIT directories: %v", err)
		panic(err)
	}
}

func LoadGitConfig(entry *GitRepo) {
	path := filepath.Join(entry.Path, ".git", "config")
	cfg, err := ini.Load(path)
	if err != nil {
		fmt.Printf("😭 failed to read %s: %v", path, err)
		os.Exit(1)
	}

	sections := cfg.Sections()
	for n := range sections {
		name := sections[n].Name()
		if strings.HasPrefix(name, "remote") {
			name = strings.Replace(name, "remote ", "", 1)
			name = strings.Replace(name, "\"", "", -1)

			remote := GitRemote{
				Name: name,
				Url:  sections[n].Key("url").String(),
			}
			entry.Remotes = append(entry.Remotes, remote)
		}
	}
}

//	func FetchGitStatus(config Configuration) GitRepo {
//		for _, entry := range config.repos {
//			entry.state = gitStatus(entry.path)
//		}
//	}

func Push(repo string, remote string) GitResult {
	var args [2]string
	args[0] = "push"
	args[1] = remote
	return execute(repo, args[:])
}

func Pull(repo string, remote string) GitResult {
	var args [2]string
	args[0] = "pull"
	args[1] = remote
	return execute(repo, args[:])
}

func Status(repo string) GitResult {
	var args [2]string
	args[0] = "status"
	args[1] = "--porcelain"
	return execute(repo, args[:])
}

func execute(dir string, args []string) GitResult {
	cmd := exec.Command("git", args[:]...)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("😭 failed to execute GIT command: %v", err)
	}

	return GitResult{string(output), cmd.ProcessState.ExitCode()}
}
