package metrics

import (
	"expvar"
	"net/http"
	"time"
)

// HTTPMetrics middleware collects comprehensive HTTP request metrics
func HTTPMetrics(next http.Handler) http.Handler {
	var (
		requestsInFlight        = expvar.NewInt("requests.in_flight")
		responsesSent           = expvar.NewInt("responses.sent")
		totalProcessingDuration = expvar.NewInt("processing.duration")
	)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Increment in-flight requests
		requestsInFlight.Add(1)

		// Call next handler
		next.ServeHTTP(w, r)

		// Calculate duration
		duration := time.Since(start)

		totalProcessingDuration.Add(int64(duration))

		responsesSent.Add(1)

		// Decrement in-flight requests
		requestsInFlight.Add(-1)

		// Path and method
		// path := r.URL.Path
		// method := r.Method
		// key := method + ":" + path
		// statusKey := string(rune(rw.statusCode/100)) + "xx"

		// Request counters
		// metricsData.requestsTotal.Add(key, 1)
		// metricsData.requestsByStatus.Add(statusKey, 1)
		// metricsData.requestsByStatus.Add(http.StatusText(rw.statusCode), 1)
		// metricsData.requestsByMethod.Add(method, 1)
		// metricsData.requestsByPath.Add(path, 1)

		// // Duration tracking
		// metricsData.durationsMu.Lock()
		// metricsData.durations = append(metricsData.durations, duration)
		// metricsData.durationsByPath[path] = append(metricsData.durationsByPath[path], duration)
		// metricsData.durationsMu.Unlock()

		// // Response size
		// metricsData.responseSizeTotal.Add(int64(rw.size))
		// metricsData.responseSizeByPath.Add(path, int64(rw.size))

		// // Track status codes
		// metricsData.statusCodes.Add(statusKey, 1)
	})
}
