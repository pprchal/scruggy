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

		action := ""
		if len(r.Form["submit"]) == 0 {
			action = r.Form["action"][0]
		} else {
			action = r.FormValue("submit")
		}

		switch action {
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
	for _, repo := range config.repos {
		html += "<tr>"

		// open terminal window
		repoInput := input("repo", repo.path)
		repoButton := button("RepoOpenTerm", "💻 "+repo.path)
		termForm := fmt.Sprintf("<form method=\"post\">%s %s</form>\n", repoInput, repoButton)
		html += "<td>" + termForm + "</td>\n"
		html += "<td>&nbsp;</td>\n"
		html += "</tr>\n"
	}

	return html
}

func button(action string, text string) string {
	return fmt.Sprintf("<button type=\"submit\" name=\"action\" value=\"%s\">%s</button>", action, text)
}

func input(name string, value string) string {
	return fmt.Sprintf("<input type=\"hidden\" name=\"%s\" value=\"%s\" />", name, value)
}
