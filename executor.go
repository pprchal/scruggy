package main

import (
	"log"
	"os/exec"
)

type GitResult struct {
	text   string
	status int
}

func gitPush(repo string, remote string) GitResult {
	var args [2]string
	args[0] = "push"
	args[1] = remote
	return execute(repo, args[:])
}

func gitPull(repo string, remote string) GitResult {
	var args [2]string
	args[0] = "pull"
	args[1] = remote
	return execute(repo, args[:])
}

func gitStatus(repo string, remote string) GitResult {
	var args [3]string
	args[0] = "pull"
	args[1] = remote
	args[2] = "--porcelain"
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
