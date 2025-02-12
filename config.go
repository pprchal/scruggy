package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

// TODO: rewrite to use channels (yield)
func loadCsv() []GitEntry {
	var gitEntries []GitEntry

	file, err := os.Open("repos.csv")
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		cols := strings.Split(scanner.Text(), ";")

		entry := GitEntry{
			path:    cols[1],
			branch:  cols[0],
			remotes: []GitRemote{},
		}
		gitEntries = append(gitEntries, entry)
	}

	return gitEntries
}

// merge ini+csv together
func loadConfiguration() Configuration {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("😭 failed to read config.ini: %v", err)
		os.Exit(1)
	}

	conf := Configuration{
		root:   cfg.Section("global").Key("root").String(),
		period: cfg.Section("global").Key("period").String(),
	}
	conf.entries = loadCsv()
	port, err := cfg.Section("global").Key("port").Int()
	if err != nil {
		log.Fatalf("😭 invalid port value: %v", err)
	}
	conf.port = port
	return conf
}
