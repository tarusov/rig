package locker_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"github.com/tarusov/rig/locker"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type redisContainer struct {
	testcontainers.Container
	URI string
}

func setupRedis(ctx context.Context) (*redisContainer, error) {

	req := testcontainers.ContainerRequest{
		Image:        "redis:6",
		ExposedPorts: []string{"6379/tcp"},
		WaitingFor:   wait.ForLog("* Ready to accept connections"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	mappedPort, err := container.MappedPort(ctx, "6379")
	if err != nil {
		return nil, err
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		return nil, err
	}

	uri := fmt.Sprintf("redis://%s:%s", hostIP, mappedPort.Port())

	return &redisContainer{Container: container, URI: uri}, nil
}

// redisServer for tests.
var redisServer *redisContainer

func TestMain(m *testing.M) {

	var (
		ctx = context.Background()
		err error
	)

	redisServer, err = setupRedis(ctx)
	if err != nil {
		log.Println("redis container not inited", err)
		os.Exit(1)
	}
	defer func() {
		if err = redisServer.Terminate(ctx); err != nil {
			log.Println("failed to terminate redis container", err)
		}
	}()

	os.Exit(m.Run())
}

func TestLockMultipleTries(t *testing.T) {

	if redisServer == nil {
		t.Skip("server not inited")
	}

	opts, err := redis.ParseURL(redisServer.URI)
	require.ErrorIsf(t, err, nil, "TestLockMultipleTries: unexpected parse uri: %v", err)

	lockerClient, err := locker.New(
		redis.NewClient(opts),
		locker.WithRetryCount(1),
		locker.WithRetryTimeout(100*time.Millisecond),
	)
	require.ErrorIsf(t, err, nil, "TestLockMultipleTries: unexpected locker error: %v", err)

	var (
		ctx = context.Background()
		key = uuid.NewString()
	)
	unlock, err := lockerClient.Lock(ctx, key, time.Second)
	require.ErrorIsf(t, err, nil, "TestLockMultipleTries: unexpected lock error: %v", err)
	defer func() {
		err := unlock()
		require.ErrorIsf(t, err, nil, "TestLockMultipleTries: unexpected unlock error: %v", err)
	}()

	_, err = lockerClient.Lock(ctx, key, time.Second)
	require.ErrorIsf(t, err, locker.ErrLockNotObtained, "TestLockMultipleTries: unexpected error: %v", err)
}

func TestLockTTL(t *testing.T) {

	if redisServer == nil {
		t.Skip("server not inited")
	}

	opts, err := redis.ParseURL(redisServer.URI)
	require.ErrorIsf(t, err, nil, "TestLockTTL: unexpected parse uri: %v", err)

	lockerClient, err := locker.New(
		redis.NewClient(opts),
		locker.WithRetryCount(1),
		locker.WithRetryTimeout(100*time.Millisecond),
	)
	require.ErrorIsf(t, err, nil, "TestLockTTL: unexpected locker error: %v", err)

	var (
		ctx = context.Background()
		key = uuid.NewString()
	)
	_, err = lockerClient.Lock(ctx, key, time.Second)
	require.ErrorIsf(t, err, nil, "TestLockTTL: unexpected lock error: %v", err)

	time.Sleep(time.Second)

	unlock, err := lockerClient.Lock(ctx, key, time.Second)
	require.ErrorIsf(t, err, nil, "TestLockTTL: unexpected lock error: %v", err)

	_ = unlock()

	// Check if not obtained.
	err = unlock()
	require.ErrorIsf(t, err, locker.ErrNotLocked, "TestLockTTL: unexpected lock error: %v", err)
}
