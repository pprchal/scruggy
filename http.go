package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func loadTemplate() string {
	path := filepath.Join("html", "index.html")
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

		switch r.Form["action"][0] {
		case "AddRepo":
			AddRepo(r.FormValue("repo"))

		case "ScanStart":
			ScanStart()

		case "ScanStop":
			ScanStop()

		case "SyncAll":
			SyncAll()

		case "RepoOpenTerm":
			OpenTerminalWindow(r.FormValue("repo"))

		case "RepoAction":
			RepoAction(r.FormValue("repo"), r.FormValue("action"), r.FormValue("remote"))
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

func StartHttp(config Configuration) {
	http.HandleFunc("/", handler)
	port := ":" + strconv.Itoa(config.port)
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}

func renderNewRepos() string {
	html := ""
	for _, repo := range config.new_repos {
		html += "<tr>"
		html += "<td>" + repo + "</td>\r\n"

		repo_input := fmt.Sprintf("<input type=\"hidden\" name=\"repo\" value=\"%s\" />", repo)
		html += fmt.Sprintf("<td><form method=\"post\" action=\"AddRepo\">%s<input type=\"submit\" value=\"➕ Add repo\" /></form></td>", repo_input)
		html += "</tr>\r\n"
	}

	return html
}

func renderRepos() string {
	html := ""
	for _, repo := range config.repos {
		html += "<tr>"

		// open terminal window
		repoInput := fmt.Sprintf("<input type=\"hidden\" name=\"repo\" value=\"%s\" />", repo.path)
		repoButton := fmt.Sprintf(
			"<input type=\"hidden\" name=\"action\" value=\"%s\" />"+
				"<input type=\"submit\" value=\"💻 %s\" />", "RepoOpenTerm", repo.path)
		repoForm := fmt.Sprintf("<form method=\"post\">%s %s</form>", repoInput, repoButton)
		html += "<td>" + repoForm + "</td>\r\n"

		// action buttons
		htmlActions := ""
		for _, action := range repo.actions {
			htmlActions += fmt.Sprintf("<form method=\"post\" action=\"RepoAction\"><input type=\"hidden\" name=\"repo\" value=\"%s\" /><input type=\"hidden\" name=\"action\" value=\"%s\" /><input type=\"hidden\" name=\"remote\" value=\"%s\" /><input type=\"submit\" value=\"⇧ %s\" /></form>", repo.path, action.action, action.remote, action.action)
		}
		html += "<td>" + htmlActions + "</td>\r\n"

		html += "</tr>\r\n"
	}

	return html
}
