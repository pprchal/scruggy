package http

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"scruggy/actions"
	"scruggy/config"
	"scruggy/git"
	"strconv"
	"strings"
)

func loadTemplate() string {
	path := filepath.Join("http", "index.html")
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Print(err)
		panic("cannot locate index.html")
	}
	return string(data)
}

func handler(writer http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			log.Print(err)
			http.Error(writer, "Unable to parse form", http.StatusBadRequest)
			return
		}

		action := ""
		if len(r.Form["submit"]) == 0 {
			action = r.Form["action"][0]
		} else {
			action = r.FormValue("submit")
		}

		switch action {
		case "AddRepo":
			actions.AddRepo(r.FormValue("repo"))

		case "ScanStart":
			actions.ScanStart()

		case "ScanStop":
			actions.ScanStop()

		case "SyncAll":
			actions.SyncAll()

		case "RepoOpenTerm":
			actions.OpenTerminalWindow(r.FormValue("repo"))

		case "RepoAction":
			actions.RepoAction(r.FormValue("repo"), r.FormValue("gitAction"), r.FormValue("remote"))

		case "Quit":
			actions.Quit()
		}
	} else if r.Method == "GET" && strings.Contains(r.RequestURI, "css") {
		path := strings.Replace(r.RequestURI, "/css/", "", 1)
		path = filepath.Join("html", "css", path)
		http.ServeFile(writer, r, path)
		return
	}

	html := strings.Replace(loadTemplate(), "{new_repos}", renderNewRepos(), -1)
	html = strings.Replace(html, "{repos}", renderRepos(), -1)
	fmt.Fprintf(writer, html)
}

func StartHttp() {
	http.HandleFunc("/", handler)
	port := ":" + strconv.Itoa(config.GlobalConfig.Port)
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}

func renderNewRepos() string {
	html := ""
	for _, repo := range config.GlobalConfig.NewRepos {
		html += "<tr>"
		html += "<td>" + repo + "</td>\r\n"

		// add repo to config
		repoInput := input("repo", repo)
		repoButton := button("AddRepo", "➕ "+repo)
		repoForm := fmt.Sprintf("<form method=\"post\">%s %s</form>\n", repoInput, repoButton)
		html += "<td>" + repoForm + "</td>\n"
		html += "</tr>\r\n"
	}

	return html
}

func renderRepos() string {
	html := ""
	for _, repo := range config.GlobalConfig.Repos {
		html += "<tr>"

		// open terminal window
		repoInput := input("repo", repo.Path)
		repoButton := button("RepoOpenTerm", "💻 "+repo.Path)
		termForm := fmt.Sprintf("<form method=\"post\">%s %s</form>\n", repoInput, repoButton)

		// sync actions ⇓ ⇑  ⇕
		actionForm := ""
		for _, action := range repo.Actions {
			actionForm += gitAction(action, repo)
		}

		html += "<td>" + termForm + "</td>\n"
		html += "<td>" + actionForm + "</td>\n"
		html += "</tr>\n"
	}

	return html
}

func gitAction(action git.GitAction, repo git.GitRepo) string {
	actionRepo := input("repo", repo.Path)
	actionRemote := input("remote", action.Remote)
	gitAction := input("gitAction", action.Action)

	symbol := "⇑"
	if action.Action == "pull" {
		symbol = "⇓"
	}

	// find remote
	remote := git.GitRemote{}
	for n := range repo.Remotes {
		if repo.Remotes[n].Name == action.Remote {
			remote = repo.Remotes[n]
			break
		}
	}

	actionButton := button("RepoAction", fmt.Sprintf("%s %s [%s]", symbol, action.Remote, remote.Url))
	actionForm := fmt.Sprintf("<form method=\"post\">%s %s %s %s</form>\n", actionRepo, actionRemote, actionButton, gitAction)
	return actionForm
}

func button(action string, text string) string {
	return fmt.Sprintf("<button type=\"submit\" name=\"action\" value=\"%s\">%s</button>", action, text)
}

func input(name string, value string) string {
	return fmt.Sprintf("<input type=\"hidden\" name=\"%s\" value=\"%s\" />", name, value)
}
