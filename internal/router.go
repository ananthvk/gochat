package internal

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/ananthvk/gochat/internal/app"
	"github.com/ananthvk/gochat/internal/helpers"
	"github.com/ananthvk/gochat/internal/realtime"
	"github.com/go-chi/chi/v5"
)

func Routes(app *app.App) chi.Router {
	router := chi.NewRouter()
	router.Mount("/realtime", realtime.Routes(app.RealtimeService))
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) { healthCheckHandler(app, w, r) })
	return router
}

func healthCheckHandler(app *app.App, w http.ResponseWriter, _ *http.Request) {
	if !app.Config.EnableDetailedHealthCheck {
		helpers.RespondWithJSON(w, 200, map[string]any{
			"status": "ok",
		})
		return
	}
	platformName := fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	memStatsMap := map[string]any{
		"alloc":           toMiB(memStats.Alloc),
		"total_alloc":     toMiB(memStats.TotalAlloc),
		"heap_alloc":      toMiB(memStats.HeapAlloc),
		"heap_sys":        toMiB(memStats.HeapSys),
		"stack_inuse":     toMiB(memStats.StackInuse),
		"memory_obtained": toMiB(memStats.Sys),
		"num_gc":          memStats.NumGC,
	}

	helpers.RespondWithJSON(w, 200, map[string]any{
		"environment": app.Config.Env,
		"platform":    platformName,
		"status":      "ok",
		"timestamp":   time.Now().UTC(),
		"uptime":      time.Since(app.StartTime).String(),
		"go_version":  runtime.Version(),
		"version":     app.Version,
		"memory":      memStatsMap,
	})
}

func toMiB(value uint64) string {
	valueInMiB := float64(value) / (1024 * 1024)
	return fmt.Sprintf("%.2f MiB", valueInMiB)
}
