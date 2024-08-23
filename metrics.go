package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var (
	resolveMissingTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "iae_image_not_found_errors_total",
			Help: "The number of times an image was not found in the registry",
		},
		[]string{"image"},
	)
)

func healthcheckHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, err := fmt.Fprintf(w, "OK")
	if err != nil {
		log.Errorf("error responding to request %v", err)
	}
}

func startMetricsServer(port int) error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.HandleFunc("/health", healthcheckHandler)
	log.Printf("Starting metrics server on port %d", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}
