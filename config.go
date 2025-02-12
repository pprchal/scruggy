package main

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// load config.ini
func loadIni() Configuration {
	conf := Configuration{root: "~/"}
	return conf
}

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
	conf := Configuration{root: "~/"}
	conf.entries = loadCsv()
	conf.port = 8080
	return conf
}
