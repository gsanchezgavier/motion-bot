package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

const namespace = "motionbot"

var (
	MovementDetectionEventsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "movement_detection_events_total",
			Help:      "Number of motion events",
		},
		[]string{"type"},
	)
	UnauthorizedRequestsTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "unauthorized_requests_total",
			Help:      "Number of unauthorized requests",
		},
	)
	HealthCheckRequestedTotal = promauto.NewCounter(
		prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "health_check_requested_total",
			Help:      "Number of health check requested",
		},
	)
	SubscribedChats = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "subscribed_chats_count",
			Help:      "Number of subscribed chats",
		},
	)
	MovementDetectionActivated = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "movement_detection_activated",
			Help:      "Movement detection is activated when this metric is not zero",
		},
	)
)

// Start creates a http server on the port selected and start serving prometheus metrics in the /metrics path
func Start(port string) {
	prom_server := http.NewServeMux()
	prom_server.Handle("/metrics", promhttp.Handler())
	go func() {
		log.Fatal().Err(http.ListenAndServe(":"+port, prom_server))
	}()
}
