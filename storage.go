package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// TODO: rewrite to use channels (yield)
func scanForGitDirs(root string) <-chan string {
	ch := make(chan string)
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() && info.Name() == ".git" {
			log.Println(path)
			ch <- strings.ReplaceAll(path, "\\.git", "")
			// gitPaths = append(gitPaths, p)
		}
		return nil
	})

	if err != nil {
		log.Fatalf("Error scanning for GIT directories: %v", err)
		panic(err)
	}

	return ch
}

func saveToCsv(gitEntries []GitEntry, output string) error {
	file, err := os.Create(output)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, entry := range gitEntries {
		line := fmt.Sprintf("%s;%s;[%s]", entry.branch, entry.path, joinRemotes(entry))
		if err := writer.Write([]string{line}); err != nil {
			return err
		}
	}
	return nil
}

func fetchBranches(entry GitEntry) {
	output := executeGit("branch", entry.path)
	println(output)
	// TODO: entry.branch = parseBranch(output)
}

func joinRemotes(entry GitEntry) string {
	return "TODO:remote"
}

func loadGitRepositories(config Configuration) <-chan GitEntry {
	ch := make(chan GitEntry)
	for _, entry := range config.entries {
		entry.state = executeGit("status", entry.path)
		ch <- entry
	}

	return ch
}
