package main

import (
	"scruggy/actions"
	"scruggy/http"
)

func main() {
	// fetch status of all repos
	actions.Refresh()

	// start gui
	http.StartHttp()
}
