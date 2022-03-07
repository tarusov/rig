package exec_test

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/oklog/run"
	"github.com/stretchr/testify/require"
	"github.com/tarusov/rig/exec"
)

func TestHealthEndpointHandler(t *testing.T) {

	var (
		ctx, cancel = context.WithCancel(context.Background())
		g           = run.Group{}
	)

	AddTestInterrupter(ctx, &g)
	exec.AddHealthEndpoint(ctx, &g, "/health", 35000)

	go func() {
		err := g.Run()
		require.ErrorIsf(t, err, ErrTestTerminated, "TestHeathEndpointHandler: group termination unexpected error: %v", err)
	}()

	httpReq, err := http.NewRequest(http.MethodGet, "http://localhost:35000/health", nil)
	require.ErrorIsf(t, err, nil, "TestHeathEndpointHandler: create request unexpected error: %v", err)

	httpResp, err := http.DefaultClient.Do(httpReq)
	require.ErrorIsf(t, err, nil, "TestHeathEndpointHandler: send request unexpected error: %v", err)

	respBody, err := ioutil.ReadAll(httpResp.Body)
	require.ErrorIsf(t, err, nil, "TestHeathEndpointHandler: read response body unexpected error: %v", err)

	require.Equal(t, string(respBody), "OK", "TestHeathEndpointHandler: unexpected response body: %s", string(respBody))

	cancel()
}

// ErrTestTerminated is custom error for group interupt.
var ErrTestTerminated = errors.New("test terminated")

// AddTestInterrupter throw "test terminated" error to group on cancel.
func AddTestInterrupter(ctx context.Context, g *run.Group) {
	g.Add(func() error {
		<-ctx.Done()
		return ErrTestTerminated
	}, func(err error) {
		// Do nothing.
	})
}
