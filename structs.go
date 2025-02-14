package main

// single remote entry
type GitRemote struct {
	name string
	url  string
}

// git entry
type GitEntry struct {
	path         string
	sync_remotes []GitRemote
	state        string
}

// man program configuration
type Configuration struct {
	root      string
	entries   []GitEntry
	port      int
	period    string
	new_repos []string
}
