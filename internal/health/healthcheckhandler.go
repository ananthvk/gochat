package health

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"time"

	"github.com/ananthvk/gochat/internal/app"
	"github.com/ananthvk/gochat/internal/helpers"
)

func HealthCheckHandler(app *app.App, w http.ResponseWriter, _ *http.Request) {
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

	ctx, cancel := context.WithTimeout(app.Ctx, app.Config.DbPingTimeout)
	defer cancel()

	start := time.Now()
	err := app.DatabaseService.Pool.Ping(ctx)
	duration := time.Since(start)

	var dbStatsMap map[string]any
	if err != nil {
		dbStatsMap = map[string]any{
			"status": "error",
		}
		slog.Error("database ping failed", "error", err)
	} else {
		dbStatsMap = map[string]any{
			"status":        "ok",
			"ping_duration": duration.String(),
		}
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
		"database":    dbStatsMap,
	})
}

func toMiB(value uint64) string {
	valueInMiB := float64(value) / (1024 * 1024)
	return fmt.Sprintf("%.2f MiB", valueInMiB)
}
