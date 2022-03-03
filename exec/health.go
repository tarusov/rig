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
func handleHealthEndpoint(w http.ResponseWriter, _ *http.Request) {
	if _, err := fmt.Fprint(w, "OK"); err != nil {
		logger.Global().WithErr(err).Error("failed to write health response")
	}
}

// AddHealthEndpoint setup health enpoint.
func AddHealthEndpoint(g *run.Group, enpoint string, healthPort int) {

	var mux = chi.NewMux()
	mux.Handle(enpoint, http.HandlerFunc(handleHealthEndpoint))

	var server = http.Server{
		Addr:    fmt.Sprintf(":%d", healthPort),
		Handler: mux,
	}

	g.Add(func() error {
		if err := server.ListenAndServe(); err != nil {
			logger.Global().WithErr(err).Error("health hanlder listen and serve error")
			return err
		}
		logger.Global().Info("health hanlder stopped")
		return nil
	}, func(error) {
		if err := server.Shutdown(context.Background()); err != nil {
			logger.Global().WithErr(err).Error("health hanlder shutdown error")
		}
		logger.Global().Info("health hanlder interrupted")
	})
}
