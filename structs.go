package main

// git remote repo - .git/config
type GitRemote struct {
	name string
	url  string
}

// config-action pair - .config.ini
type GitAction struct {
	remote string
	action string
}

// git repository - .config.ini
type GitRepo struct {
	path    string
	remotes []GitRemote
	actions []GitAction
	state   string
}

// main program configuration - global state
type Configuration struct {
	root      string
	repos     []GitRepo
	port      int
	period    string
	new_repos []string
}
