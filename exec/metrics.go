package exec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oklog/run"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/tarusov/rig/logger"
	"github.com/tarusov/rig/metrics"
)

// AddMetricsEndpoint setup prometheus metrics handler.
func AddMetricsEndpoint(ctx context.Context, g *run.Group, registry metrics.Registry, enpoint string, metricsPort int) {

	var mux = chi.NewMux()
	mux.Handle(enpoint, promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	var server = http.Server{
		Addr:    fmt.Sprintf(":%d", metricsPort),
		Handler: mux,
	}

	g.Add(func() error {

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.FromContext(ctx).WithErr(err).Error("prometheus metrics hanlder listen and serve error")
			return err
		}

		logger.FromContext(ctx).Info("prometheus metrics hanlder stopped")
		return nil

	}, func(error) {

		if err := server.Shutdown(context.Background()); err != nil {
			logger.FromContext(ctx).WithErr(err).Error("prometheus metrics hanlder shutdown error")
			return
		}

		logger.FromContext(ctx).Info("prometheus metrics hanlder interrupted")
	})
}
