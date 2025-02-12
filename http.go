package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

func loadHtml() string {
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
		log.Print("POST request received")

		request, err := io.ReadAll(r.Body)
		if err != nil {
			log.Print(err)
		}
		log.Print(string(request))
	}

	html := loadHtml()
	fmt.Fprintf(writer, html)
}

func startHttp(config Configuration) {
	http.HandleFunc("/", handler)
	port := ":" + strconv.Itoa(config.port)
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}
