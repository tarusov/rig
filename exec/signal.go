package exec

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/oklog/run"
	"github.com/tarusov/rig/logger"
)

// AddSignalWatcher setup a signal recivier for run group.
func AddSignalWatcher(ctx context.Context, g *run.Group) {

	var sCtx, sCancel = context.WithCancel(ctx)

	g.Add(func() error {

		sigChan := make(chan os.Signal, 2)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		logger.FromContext(sCtx).Info("signal watcher started")

		select {
		case c := <-sigChan:
			return fmt.Errorf("terminated with sig %q", c)
		case <-sCtx.Done():
			return nil
		}

	}, func(err error) {
		sCancel()
	})
}
