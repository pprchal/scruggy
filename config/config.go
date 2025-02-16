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
	for n := range sections {
		section := sections[n]
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
	actionSplits := strings.Split(remotes, ",")
	actions := []git.GitAction{}
	for i := range actionSplits {
		action := git.GitAction{
			Action: "",
		}

		if strings.HasPrefix(actionSplits[i], "push-") {
			action.Action = "push"
			action.Remote = strings.TrimPrefix(actionSplits[i], "push-")
		} else if strings.HasPrefix(actionSplits[i], "pull-") {
			action.Action = "pull"
			action.Remote = strings.TrimPrefix(actionSplits[i], "pull-")
		} else {
			panic("😭 invalid action: " + actionSplits[i])
		}

		actions = append(actions, action)
	}

	return actions
}
