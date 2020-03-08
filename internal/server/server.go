package server

import (
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/telliott-io/front/pkg/projects"
)

func New(
	loader projects.Loader,
	env string,
) (http.Handler, error) {
	return &server{
		loader: loader,
		env:    env,
	}, nil
}

type server struct {
	loader projects.Loader
	env    string
}

var (
	indexServed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "front_index_served",
		Help: "The total number index requests served",
	})
)

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	indexServed.Inc()

	tracer := opentracing.GlobalTracer()
	span := tracer.StartSpan("serve_index")

	span.SetTag("path", r.URL.Path)
	span.SetTag("url", r.URL.String())
	span.SetTag("referer", r.Referer())
	span.SetTag("user-agent", r.UserAgent())

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
		Env   string
	}{
		Items: p,
		Env:   s.env,
	}

	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	renderPageSpan.Finish()
}
