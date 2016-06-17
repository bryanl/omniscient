package omniscient

import (
	"fmt"
	"time"

	"github.com/labstack/echo"
	"github.com/prometheus/client_golang/prometheus"
)

var hitCounter = prometheus.NewCounter(prometheus.CounterOpts{
	Name: "omniscient_service_hits",
	Help: "Number of service hits.",
})

var requestHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
	Name:    "omniscient_request_duration_hist",
	Help:    "Request duration.",
	Buckets: prometheus.LinearBuckets(20, 5, 5),
})

var requestSummary = prometheus.NewSummary(prometheus.SummaryOpts{
	Name: "omniscient_request_duration_sum",
	Help: "Request duration.",
})

var metricsCollectors = []prometheus.Collector{
	hitCounter,
	requestHistogram,
	requestSummary,
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
		return func(c echo.Context) error {
			if err := h(c); err != nil {
				c.Error(err)
			}

			start := time.Now()
			hitCounter.Inc()

			dur := time.Since(start)
			fmt.Println("dur", dur.Seconds())

			requestHistogram.Observe(dur.Seconds())
			requestSummary.Observe(dur.Seconds())

			return nil
		}
	}
}
