package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

// merge ini+csv together
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

	// [/repos...]
	conf.entries = LoadGitEntries(cfg)
	return conf
}

func LoadGitEntries(cfg *ini.File) []GitEntry {
	gitEntries := []GitEntry{}
	sections := cfg.Sections()
	for n := range sections {
		section := sections[n]
		if strings.HasPrefix(section.Name(), "global") {
			continue
		}

		if strings.HasPrefix(section.Name(), "DEFAULT") {
			continue
		}

		gitEntries = append(gitEntries, GitEntry{
			path:         section.Name(),
			sync_remotes: ParseRemotes(section.KeysHash()["sync_remotes"]),
		})
	}

	return gitEntries
}

func ParseRemotes(remotes string) []GitRemote {
	sync_remotes := strings.Split(remotes, ",")
	gitRemotes := []GitRemote{}
	for i := range sync_remotes {
		gitRemotes = append(gitRemotes, GitRemote{
			name: sync_remotes[i],
			url:  "",
		})
	}

	return gitRemotes
}
