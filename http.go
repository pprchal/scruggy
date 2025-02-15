package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func loadTemplate() string {
	data, err := os.ReadFile("index.html")
	if err != nil {
		fmt.Print(err)
		panic("cannot locate index.html")
	}
	html := string(data)
	return html
}

func handler(writer http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		if err := r.ParseForm(); err != nil {
			log.Print(err)
			http.Error(writer, "Unable to parse form", http.StatusBadRequest)
			return
		}

		switch r.RequestURI {
		case "/AddRepo":
			AddRepo(r.FormValue("repo"))

		case "/ScanStart":
			ScanStart()

		case "/ScanStop":
			ScanStop()

		case "/SyncAll":
			SyncAll()

		case "/RepoAction":
			RepoAction(r.FormValue("repo"), r.FormValue("action"), r.FormValue("remote"))
		}
	}

	// for key, values := range r.Form {
	// 	for _, value := range values {
	// 		log.Printf("Form field %s: %s", key, value)
	// 	}
	// }

	html := strings.Replace(loadTemplate(), "{new_repos}", renderNewRepos(), -1)
	html = strings.Replace(html, "{repos}", renderRepos(), -1)
	fmt.Fprintf(writer, html)
}

func startHttp(config Configuration) {
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
		html += "<td>" + repo.path + "</td>\r\n"

		htmlActions := ""
		for _, action := range repo.actions {
			htmlActions += fmt.Sprintf("<form method=\"post\" action=\"RepoAction\"><input type=\"hidden\" name=\"repo\" value=\"%s\" /><input type=\"hidden\" name=\"action\" value=\"%s\" /><input type=\"hidden\" name=\"remote\" value=\"%s\" /><input type=\"submit\" value=\"%s\" /></form>", repo.path, action.action, action.remote, action.action)
		}
		html += "<td>" + htmlActions + "</td>\r\n"

		// repo_input := fmt.Sprintf("<input type=\"hidden\" name=\"repo\" value=\"%s\" /><input type=\"hidden\" name=\"action\" value=\"%s\" />", repo, "pull")
		// html += fmt.Sprintf("<td><form method=\"post\" action=\"RepoPush\">%s<input type=\"submit\" value=\"sync\" /></form></td>", repo_input)
		html += "</tr>\r\n"
	}

	return html
}
