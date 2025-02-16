package main

import (
	"log"
	"os/exec"
)

func ExecuteGit(command string, dir string) string {
	// --porcelain
	cmd := exec.Command("git", command)
	cmd.Dir = dir
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("😭 failed to execute GIT command %s: %v", command, err)
	}

	return string(output)
}
