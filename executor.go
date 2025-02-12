package main

import (
	"log"
	"os/exec"
)

func executeGit(command string, dir string) string {
	cmd := exec.Command("git", command)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("😭 failed to execute GIT command: %v", err)
	}

	return string(output)
}
