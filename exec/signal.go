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
func AddSignalWatcher(g *run.Group) {

	var ctx, cancel = context.WithCancel(context.Background())

	g.Add(func() error {
		sigChan := make(chan os.Signal, 2)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		logger.Global().Info("signal watcher started")
		select {
		case c := <-sigChan:
			return fmt.Errorf("terminated with sig %q", c)
		case <-ctx.Done():
			return nil
		}
	}, func(err error) {
		cancel()
	})
}
