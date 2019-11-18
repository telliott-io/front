package main

import (
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("public/styles/"))
	http.Handle("/styles/", http.StripPrefix("/styles/", fs))
	views := http.FileServer(http.Dir("views/"))
	http.Handle("/", views)

	log.Fatal(http.ListenAndServe(":80", nil))
}
