package omniscient

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"gopkg.in/labstack/echo.v1"
)

var hitCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "omniscient_service_hits",
	Help: "Number of service hits.",
})

var requestHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name:    "omniscient_request_duration",
	Help:    "Request duration.",
	Buckets: prometheus.LinearBuckets(20, 5, 5),
})

var metricsCollectors = []prometheus.Collector{
	hitCounter,
	requestHistogram,
}

func initMetrics() error {
	for _, c := range metricsCollectors {
		err := prometheus.Register(c)
		if err != nil {
			return err
		}

	}
	return nil
}

// HitCounter tracks hits.
func HitCounter() echo.MiddlewareFunc {
	return func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c *echo.Context) error {
			if err := h(c); err != nil {
				c.Error(err)
			}

			before := time.Now()
			hitCounter.Inc()
			now := time.Now()
			dur := now.Sub(before)

			requestHistogram.Observe(float64(dur))

			return nil
		}
	}
}
