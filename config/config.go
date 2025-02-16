package config

import (
	"fmt"
	"log"
	"os"
	"scruggy/git"
	"strings"

	"gopkg.in/ini.v1"
)

// main program configuration - global state
type Configuration struct {
	Root     string
	Repos    []git.GitRepo
	Port     int
	NewRepos []string
}

var GlobalConfig Configuration

func LoadConfiguration() Configuration {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("😭 failed to read config.ini: %v", err)
		os.Exit(1)
	}

	// [global]
	conf := Configuration{
		Root: cfg.Section("global").Key("scan").String(),
	}

	port, err := cfg.Section("global").Key("port").Int()
	if err != nil {
		log.Fatalf("😭 invalid port value[%s]: %v", cfg.Section("global").Key("port"), err)
	}
	conf.Port = port

	// load .git/config
	conf.Repos = BuildGitRepos(cfg)
	return conf
}

func BuildGitRepos(cfg *ini.File) []git.GitRepo {
	repos := []git.GitRepo{}
	sections := cfg.Sections()
	for _, section := range sections {
		if strings.HasPrefix(section.Name(), "global") {
			continue
		}

		if strings.HasPrefix(section.Name(), "DEFAULT") {
			continue
		}

		repo := git.GitRepo{
			Path:    section.Name(),
			Actions: ParseActions(section.KeysHash()["actions"]),
			State:   "",
		}

		git.LoadGitConfig(&repo)
		repos = append(repos, repo)
	}

	return repos
}

func ParseActions(remotes string) []git.GitAction {
	actions := []git.GitAction{}
	for _, actionSplit := range strings.Split(remotes, ",") {
		if actionSplit == "" {
			continue
		}

		action := git.GitAction{
			Action: "",
		}

		if strings.HasPrefix(actionSplit, "push-") {
			action.Action = "push"
			action.Remote = strings.TrimPrefix(actionSplit, "push-")
		} else if strings.HasPrefix(actionSplit, "pull-") {
			action.Action = "pull"
			action.Remote = strings.TrimPrefix(actionSplit, "pull-")
		} else {
			panic("😭 invalid action: " + actionSplit)
		}

		actions = append(actions, action)
	}

	return actions
}
