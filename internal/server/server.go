package server

import (
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/telliott-io/front/pkg/projects"
)

func New(
	loader projects.Loader,
) (http.Handler, error) {
	return &server{
		loader: loader,
	}, nil
}

type server struct {
	loader projects.Loader
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tpl, err := ioutil.ReadFile("views/index.html")
	if err != nil {
		// TODO: Handle file not existing
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	t, err := template.New("page").Parse(string(tpl))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	p, err := s.loader.GetProjects()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		Items []projects.Project
	}{
		Items: p,
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
