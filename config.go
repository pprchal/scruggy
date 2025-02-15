package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

func LoadConfiguration() Configuration {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("😭 failed to read config.ini: %v", err)
		os.Exit(1)
	}

	// [global]
	conf := Configuration{
		root:   cfg.Section("global").Key("scan").String(),
		period: cfg.Section("global").Key("period").String(),
	}

	port, err := cfg.Section("global").Key("port").Int()
	if err != nil {
		log.Fatalf("😭 invalid port value[%s]: %v", cfg.Section("global").Key("port"), err)
	}
	conf.port = port

	// load .git/config
	conf.repos = BuildGitRepos(cfg)
	return conf
}

func BuildGitRepos(cfg *ini.File) []GitRepo {
	repos := []GitRepo{}
	sections := cfg.Sections()
	for n := range sections {
		section := sections[n]
		if strings.HasPrefix(section.Name(), "global") {
			continue
		}

		if strings.HasPrefix(section.Name(), "DEFAULT") {
			continue
		}

		repo := GitRepo{
			path:    section.Name(),
			actions: ParseActions(section.KeysHash()["actions"]),
			state:   "",
		}

		LoadGitConfig(&repo)
		repos = append(repos, repo)
	}

	return repos
}

func ParseActions(remotes string) []GitAction {
	actionSplits := strings.Split(remotes, ",")
	actions := []GitAction{}
	for i := range actionSplits {
		action := GitAction{
			action: "",
		}

		if strings.HasPrefix(actionSplits[i], "push-") {
			action.action = "push"
			action.remote = strings.TrimPrefix(actionSplits[i], "push-")
		} else if strings.HasPrefix(actionSplits[i], "pull-") {
			action.action = "pull"
			action.remote = strings.TrimPrefix(actionSplits[i], "pull-")
		} else {
			panic("😭 invalid action: " + actionSplits[i])
		}

		actions = append(actions, action)
	}

	return actions
}
