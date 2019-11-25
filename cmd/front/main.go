package main

import (
	"io"
	"log"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"

	"github.com/telliott-io/front/internal/server"
	"github.com/telliott-io/front/pkg/projects/cachingloader"
	"github.com/telliott-io/front/pkg/projects/kubernetesloader"
)

func main() {
	closer, err := setupOpenTracing()
	if err != nil {
		log.Println("Opentracing setup failed: ", err)
	}
	defer closer.Close()

	setupStaticServing()
	setupDynamicServing()

	log.Println("Listening on port 80")
	log.Fatal(http.ListenAndServe(":80", nil))
}

func setupOpenTracing() (io.Closer, error) {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: "jaeger-agent:6831",
			LogSpans:           true,
		},
	}
	tracer, closer, err := cfg.New("telliott-io/front", config.Logger(jaeger.StdLogger))
	if err != nil {
		return nil, err
	}

	opentracing.SetGlobalTracer(tracer)
	return closer, nil
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

	s, err := server.New(cachingloader.New(loader))
	if err != nil {
		panic(err)
	}

	http.Handle("/", s)
}
