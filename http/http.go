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

func StartHttp() {
	http.HandleFunc("/", handler)
	port := ":" + strconv.Itoa(config.GlobalConfig.Port)
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}

func loadTemplate(name string) string {
	path := filepath.Join("http", name)
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Print(err)
		log.Panicf("😭 cannot load: %s", name)
	}
	return string(data)
}

func handler(writer http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			log.Print(err)
			http.Error(writer, "😭 unable to parse form", http.StatusBadRequest)
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
			actions.RepoActions(r.FormValue("repo"), r.FormValue("actions"))

		case "Quit":
			actions.Quit()
		}
	} else if r.Method == "GET" && strings.Contains(r.RequestURI, "css") {
		path := strings.Replace(r.RequestURI, "/css/", "", 1)
		path = filepath.Join("html", "css", path)
		http.ServeFile(writer, r, path)
		return
	}

	html := loadTemplate("index.html")
	html = strings.Replace(html, "{renderNewRepos}", renderNewRepos(), -1)
	html = strings.Replace(html, "{renderRepos}", renderRepos(), -1)
	fmt.Fprintf(writer, html)
}

func renderNewRepos() string {
	html := ""
	for _, repo := range config.GlobalConfig.NewRepos {
		html += "<tr>"
		html += "<td>" + repo + "</td>\r\n"

		// add repo to config
		repoForm := fmt.Sprintf("<form method=\"post\">%s %s</form>\n",
			input("repo", repo),
			button("AddRepo", "➕ "+repo))
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
		html += "<td>" + termForm(repo) + "</td>\n"

		// sync actions ⇓ ⇑  ⇕
		html += "<td>" + actionForm(repo) + "</td>\n"

		html += "</tr>\n"
	}

	return html
}

func termForm(repo git.GitRepo) string {
	return fmt.Sprintf("<form method=\"post\">%s %s</form>\n",
		input("repo", repo.Path),
		button("RepoOpenTerm", "💻 "+repo.Path))
}

func actionForm(repo git.GitRepo) string {
	symbols := ""
	actions := ""
	remotes := ""

	for _, action := range repo.Actions {
		switch action.Action {
		case "pull":
			symbols += "⇓"

		case "push":
			symbols += "⇑"
		}

		if remotes == "" {
			remote := findRemote(repo, action.Remote)
			remotes = fmt.Sprintf("%s [%s]", remote.Name, remote.Url)
		}
		actions += fmt.Sprintf("%s-%s,", action.Action, action.Remote)
	}

	return fmt.Sprintf("<form method=\"post\">%s %s %s</form>\n",
		input("actions", actions),
		input("repo", repo.Path),
		button("RepoAction", fmt.Sprintf("%s %s", symbols, remotes)))

}

func button(action string, text string) string {
	return fmt.Sprintf("<button type=\"submit\" name=\"action\" value=\"%s\">%s</button>", action, text)
}

func input(name string, value string) string {
	return fmt.Sprintf("<input type=\"hidden\" name=\"%s\" value=\"%s\" />", name, value)
}

func findRemote(repo git.GitRepo, name string) git.GitRemote {
	for _, remote := range repo.Remotes {
		if remote.Name == name {
			return remote
		}
	}

	return git.GitRemote{}
}
