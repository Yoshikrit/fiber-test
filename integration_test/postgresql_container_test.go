package integration_test

import (
    "context"
	"fmt"
    "database/sql"
    "testing"

    "github.com/testcontainers/testcontainers-go"
    "github.com/testcontainers/testcontainers-go/wait"
    _ "github.com/lib/pq" // Import the PostgreSQL driver
)

func TestPostgreSQLIntegration(t *testing.T) {
    ctx := context.Background()

    // Define a PostgreSQL container
    req := testcontainers.ContainerRequest{
        Image:        "postgres:latest",
        ExposedPorts: []string{"5432/tcp"},
        Env: map[string]string{
            "POSTGRES_USER":     "postgres",
            "POSTGRES_PASSWORD": "postgres",
            "POSTGRES_DB":       "postgres",
        },
        WaitingFor: wait.ForLog("database system is ready to accept connections"),
    }

    // Create the container
    container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
        ContainerRequest: req,
        Started:          true,
    })
    if err != nil {
        t.Fatal(err)
    }
	defer container.Terminate(ctx)

    // Get the PostgreSQL port
    host, err := container.Host(ctx)
    if err != nil {
        t.Fatal(err)
    }

    port, err := container.MappedPort(ctx, "5432")
    if err != nil {
        t.Fatal(err)
    }

    dsn := fmt.Sprintf("user=postgres password=postgres dbname=postgres sslmode=disable host=%s port=%s", host, port.Port())

    // Connect to the database
    db, err := sql.Open("postgres", dsn)
    if err != nil {
        t.Fatal(err)
    }
    defer db.Close()

    // Your database test code here

    // Clean up
    if err := container.Terminate(ctx); err != nil {
        t.Fatal(err)
    }
}