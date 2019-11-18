package main

import (
	"fmt"
	"html"
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("public/styles/"))
	http.Handle("/styles/", http.StripPrefix("/styles/", fs))
	views := http.FileServer(http.Dir("views/"))
	http.Handle("/", http.StripPrefix("/", views))

	http.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Println("Listening on port 80")
	log.Fatal(http.ListenAndServe(":80", nil))
}
