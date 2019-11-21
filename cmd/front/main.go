package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	fs := http.FileServer(http.Dir("public/styles/"))
	http.Handle("/styles/", http.StripPrefix("/styles/", fs))

	http.HandleFunc("/", handleIndex)

	log.Println("Listening on port 80")
	log.Fatal(http.ListenAndServe(":80", nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	tpl, err := ioutil.ReadFile("views/index.html")
	if err != nil {
		// TODO: Handle file not existing
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	t, err := template.New("page").Parse(string(tpl))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	data := struct {
		Items []struct {
			Name        string
			Description string
		}
	}{
		Items: []struct {
			Name        string
			Description string
		}{
			{
				Name:        "Item A",
				Description: "Desc",
			},
			{
				Name:        "Item B",
				Description: "Desc",
			},
			{
				Name:        "Item C",
				Description: "Desc",
			},
			{
				Name:        "Item D",
				Description: "Desc",
			},
			{
				Name:        "Item E",
				Description: "Desc",
			},
			{
				Name:        "Item F",
				Description: "Desc",
			},
		},
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
