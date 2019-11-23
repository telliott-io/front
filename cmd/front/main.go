package main

import (
	"log"
	"net/http"

	"github.com/telliott-io/front/internal/server"
	"github.com/telliott-io/front/pkg/projects/kubernetesloader"
)

func main() {

	setupStaticServing()
	setupDynamicServing()

	log.Println("Listening on port 80")
	log.Fatal(http.ListenAndServe(":80", nil))
}

func setupStaticServing() {
	fs := http.FileServer(http.Dir("public/styles/"))
	http.Handle("/styles/", http.StripPrefix("/styles/", fs))
}

func setupDynamicServing() {
	loader, err := kubernetesloader.New()
	if err != nil {
		panic(err)
	}

	s, err := server.New(loader)
	if err != nil {
		panic(err)
	}

	http.Handle("/", s)
}
