package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"gopkg.in/ini.v1"
)

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
	path := filepath.Join(entry.path, ".git", "config")
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
				name: name,
				url:  sections[n].Key("url").String(),
			}
			entry.remotes = append(entry.remotes, remote)
		}
	}
}

func FetchGitStatus(config Configuration) <-chan GitRepo {
	ch := make(chan GitRepo)

	go func() {
		defer close(ch)
		for _, entry := range config.repos {
			entry.state = ExecuteGit("status", entry.path)
			ch <- entry
		}
	}()

	return ch
}

func OpenTerminalWindow(path string) {
	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", path).Start()
	case "windows":
		exec.Command("rundll32", "url.dll,FileProtocolHandler", path).Start()
	case "darwin":
		exec.Command("open", path).Start()
	default:
		log.Fatal("unsupported platform")
	}
}
