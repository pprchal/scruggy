package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/ini.v1"
)

func FindGitRepositories(root string, f func(string)) {
	// ch := make(chan string)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			f(strings.ReplaceAll(path, "\\.git", ""))
			// ch <- strings.ReplaceAll(path, "\\.git", "")
		}
		return nil
	})

	if err != nil {
		log.Fatalf("😭 error scanning for GIT directories: %v", err)
		panic(err)
	}

	// return ch
}

func LoadGitConfig(entry GitEntry) {
	cfg, err := ini.Load(filepath.Join(entry.path, ".git", "config"))
	if err != nil {
		fmt.Printf("😭 failed to read %s/.git/config: %v", entry.path, err)
		os.Exit(1)
	}

	sections := cfg.Sections()
	for n := range sections {
		name := sections[n].Name()
		if strings.HasPrefix(name, "remote \"") {
			name = strings.TrimPrefix(name, "remote \"")
			entry.sync_remotes = append(entry.sync_remotes, GitRemote{name: name})
		}
		log.Println(name)
	}
}

func FetchGitStatus(config Configuration) <-chan GitEntry {
	ch := make(chan GitEntry)

	go func() {
		defer close(ch)
		for _, entry := range config.entries {
			entry.state = executeGit("status", entry.path)
			ch <- entry
		}
	}()

	return ch
}
