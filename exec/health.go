package exec

import (
	"context"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/oklog/run"
	"github.com/tarusov/rig/logger"
)

// handleHealthEndpoint write health response.
func handleHealthEndpoint(w http.ResponseWriter, r *http.Request) {
	if _, err := fmt.Fprint(w, "OK"); err != nil {
		logger.FromContext(r.Context()).WithErr(err).Error("failed to write health response")
	}
}

// AddHealthEndpoint setup health enpoint.
func AddHealthEndpoint(ctx context.Context, g *run.Group, enpoint string, healthPort int) {

	var mux = chi.NewMux()
	mux.Handle(enpoint, http.HandlerFunc(handleHealthEndpoint))

	var server = http.Server{
		Addr:    fmt.Sprintf(":%d", healthPort),
		Handler: mux,
	}

	g.Add(func() error {

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.FromContext(ctx).WithErr(err).Error("health hanlder listen and serve error")
			return err
		}

		logger.FromContext(ctx).Info("health hanlder stopped")
		return nil

	}, func(error) {

		if err := server.Shutdown(ctx); err != nil {
			logger.FromContext(ctx).WithErr(err).Error("health hanlder shutdown error")
			return
		}

		logger.FromContext(ctx).Info("health hanlder interrupted")
	})
}
