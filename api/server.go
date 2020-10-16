package api

import (
	"net/http"
)

func Start() {
	// Start up http server
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		//log.Println("Got request for:", r.URL.String())
	})
	go http.ListenAndServe(":8080", nil)
}
