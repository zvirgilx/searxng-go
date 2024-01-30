package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

func init() {
	prometheus.MustRegister(ResponseCounter)
	prometheus.MustRegister(EnginesResponseCounter)
	prometheus.MustRegister(EnginesSearchResultCounter)

}

var (
	// ResponseCounter monitors the QPS (query per second), status code of the API and the latency of the API.
	ResponseCounter = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "searxng_response_total",
			Help:    "Total response status of the server API",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "endpoint", "code"},
	)

	// EnginesResponseCounter monitors the QPS (query per second), status of the engine and the latency of the engine.
	EnginesResponseCounter = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "engines_response_total",
			Help: `The response status of the engine search result. Engine 
            usually has four results, ok (at least one result), empty 
            (no search result), block (request is blocked), and error (other failures).`,
			Buckets: prometheus.DefBuckets,
		},
		[]string{"engine", "status"},
	)

	// EnginesSearchResultCounter counts the number of results of the engine.
	EnginesSearchResultCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "engines_search_result_total",
			Help: "Total number of engines search result.",
		},
		[]string{"engine"},
	)
)
