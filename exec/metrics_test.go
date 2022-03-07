package exec_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/oklog/run"
	"github.com/stretchr/testify/require"
	"github.com/tarusov/rig/exec"
	"github.com/tarusov/rig/metrics"
)

func TestMetricsEndpointHandler(t *testing.T) {

	var (
		ctx, cancel = context.WithCancel(context.Background())
		g           = run.Group{}
	)

	AddTestInterrupter(ctx, &g)
	exec.AddMetricsEndpoint(ctx, &g, metrics.New(), "/metrics", 35000)

	go func() {
		err := g.Run()
		require.ErrorIsf(t, err, ErrTestTerminated, "TestMetricsEndpointHandler: group termination unexpected error: %v", err)
	}()

	httpReq, err := http.NewRequest(http.MethodGet, "http://localhost:35000/metrics", nil)
	require.ErrorIsf(t, err, nil, "TestMetricsEndpointHandler: create request unexpected error: %v", err)

	httpResp, err := http.DefaultClient.Do(httpReq)
	require.ErrorIsf(t, err, nil, "TestMetricsEndpointHandler: send request unexpected error: %v", err)

	respBody, err := ioutil.ReadAll(httpResp.Body)
	require.ErrorIsf(t, err, nil, "TestMetricsEndpointHandler: read response body unexpected error: %v", err)

	require.NotEmptyf(t, respBody, "TestMetricsEndpointHandler: unexpected response body: %s", string(respBody))

	cancel()
}
