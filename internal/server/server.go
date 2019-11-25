package server

import (
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/opentracing/opentracing-go"
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
	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("serve_index")

	ctx := opentracing.ContextWithSpan(r.Context(), span)
	defer span.Finish()

	templateSpan, _ := opentracing.StartSpanFromContext(ctx, "load-template")
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
	templateSpan.Finish()

	getProjectsSpan, _ := opentracing.StartSpanFromContext(ctx, "get-projects")
	p, err := s.loader.GetProjects(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	getProjectsSpan.Finish()

	renderPageSpan, _ := opentracing.StartSpanFromContext(ctx, "render-page")
	data := struct {
		Items []projects.Project
	}{
		Items: p,
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	renderPageSpan.Finish()
}
