package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
)

func renderPage() string {
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

		for key, values := range r.Form {
			for _, value := range values {
				log.Printf("Form field %s: %s", key, value)
			}
		}
	}

	html := renderPage()
	fmt.Fprintf(writer, html)
}

func startHttp(config Configuration) {
	http.HandleFunc("/", handler)
	port := ":" + strconv.Itoa(config.port)
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}
