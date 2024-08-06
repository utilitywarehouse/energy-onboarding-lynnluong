package store

import (
	"context"
	"testing"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// This code uses go TestContainers, for more inoformation see: https://golang.testcontainers.org/
//
// There is currently a bug in the latest version of go-testcontainers, if you see this error:
// ../../../go/pkg/mod/github.com/testcontainers/testcontainers-go@v0.17.0/docker.go:220:24: undefined: container.StopOptions
// Then add the following to your go.mod file
// replace (
//     github.com/docker/docker => github.com/docker/docker v20.10.3-0.20221013203545-33ab36d6b304+incompatible // 22.06 branch
// )
// See: https://golang.testcontainers.org/quickstart/#2-install-testcontainers-for-go

const (
	dbName     = "projection"
	pgUser     = "energy"
	pgPassword = "password"
	pgPort     = "5432"
	pgPortTCP  = pgPort + "/tcp"
)

func SetupTestDatabaseContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "postgres:13-alpine",
		ExposedPorts: []string{pgPortTCP},
		WaitingFor:   wait.ForListeningPort(pgPortTCP),
		Env: map[string]string{
			"POSTGRES_DB":       dbName,
			"POSTGRES_PASSWORD": pgPassword,
			"POSTGRES_USER":     pgUser,
		},
	}

	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

func GetPostgresDSN(ctx context.Context, container testcontainers.Container) (string, error) {
	host, err := container.Host(ctx)
	if err != nil {
		return "", err
	}
	port, err := container.MappedPort(ctx, pgPort)
	if err != nil {
		return "", err
	}
	return postgres.DSN(host, port.Int(), pgUser, pgPassword, dbName), nil
}

func TestStore_UpdateService(t *testing.T) {
	ctx := context.Background()

	container, err := SetupTestDatabaseContainer(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer container.Terminate(ctx)

	pgDSN, err := GetPostgresDSN(ctx, container)
	if err != nil {
		t.Fatal(err)
	}

	store, err := Setup(ctx, pgDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	// Start your tests here:
}
