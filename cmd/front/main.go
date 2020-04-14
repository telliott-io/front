package main

import (
	"io"
	"log"
	"net/http"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"

	"github.com/telliott-io/front/internal/server"
	"github.com/telliott-io/front/pkg/observability"
	"github.com/telliott-io/front/pkg/projects/cachingloader"
	"github.com/telliott-io/front/pkg/projects/kubernetesloader"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/hashicorp/consul/api"
)

func main() {
	closer, err := setupOpenTracing()
	if err != nil {
		log.Println("Opentracing setup failed: ", err)
	}
	if closer != nil {
		defer closer.Close()
	}

	setupMetrics()

	setupFaviconServing()
	setupStaticServing()
	setupDynamicServing()

	log.Println("Listening on port 80")
	log.Fatal(http.ListenAndServe(":80", nil))
}

func setupMetrics() {
	http.Handle("/metrics", promhttp.Handler())
}

func setupOpenTracing() (io.Closer, error) {
	cfg := &config.Configuration{
		Sampler: &config.SamplerConfig{
			Type:  "const",
			Param: 1,
		},
		Reporter: &config.ReporterConfig{
			LocalAgentHostPort: "jaeger-agent.monitoring.svc.cluster.local:6831",
		},
	}
	tracer, closer, err := cfg.New("telliott-io/front", config.Logger(jaeger.StdLogger))
	if err != nil {
		return nil, err
	}

	opentracing.SetGlobalTracer(
		observability.NewMetricsTracer(
			"front",
			tracer,
		),
	)
	return closer, nil
}

func setupStaticServing() {
	http.Handle("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("public/styles/"))))
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("public/assets/"))))
}

func setupFaviconServing() {
	http.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "not found", http.StatusNotFound)
	})
}

func setupDynamicServing() {
	loader, err := kubernetesloader.New()
	if err != nil {
		panic(err)
	}

	// Get a new Consul client
	var environmentName = "unknown"
	cfg := api.DefaultConfig()
	cfg.Address = "consul.consul.svc.cluster.local:8500"
	client, err := api.NewClient(cfg)
	if err == nil {
		// Get a handle to the KV API
		kv := client.KV()

		// Lookup the pair
		pair, _, err := kv.Get("deployment/name", nil)
		if err == nil {
			environmentName = string(pair.Value)
		} else {
			log.Println("Could not read KV: ", err)
		}
	} else {
		log.Println("Could not connect to consul: ", err)
	}
	log.Println("Environment name: ", environmentName)

	s, err := server.New(
		cachingloader.New(loader),
		environmentName,
	)
	if err != nil {
		panic(err)
	}

	http.Handle("/", s)
}
