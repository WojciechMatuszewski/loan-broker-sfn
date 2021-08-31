package main_test

import (
	"context"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestSFNContainer(t *testing.T) {
	ctx := context.Background()
	port, err := nat.NewPort("tcp", "8083")
	if err != nil {
		t.Fatalf(err.Error())
	}

	req := testcontainers.ContainerRequest{
		Image:        "amazon/aws-stepfunctions-local",
		ExposedPorts: []string{"8083/tcp"},
		WaitingFor:   wait.ForListeningPort(port),
	}

	sfnContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatalf(err.Error())
	}

	defer sfnContainer.Terminate(ctx)
}
