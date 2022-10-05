package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/docker/go-connections/nat"
	_ "github.com/jackc/pgx/v5/pgxpool"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

/*
func SetupPostgres(t *testing.T) (tc.Container, TodoRepository) {
	postgresPort := nat.Port("5432/tcp")
	postgres, err := tc.GenericContainer(context.Background(),
		tc.GenericContainerRequest{
			ContainerRequest: tc.ContainerRequest{
				Image:        "postgres:14.5",
				ExposedPorts: []string{postgresPort.Port()},
				Env: map[string]string{
					"POSTGRES_PASSWORD": "password",
					"POSTGRES_USER":     "user",
				},
				WaitingFor: wait.ForAll(
					wait.ForLog("database system is ready to accept connections"),
					wait.ForListeningPort(postgresPort),
				),
			},
			Started: true,
		})
	if err != nil {
		t.Fatal(err)
		return nil, nil
	}

	hostPort, err := postgres.MappedPort(context.Background(), postgresPort)
	if err != nil {
		t.Fatal(err)
		return nil, nil
	}

	postgresUrl := fmt.Sprintf("postgres://user:password@localhost:%s?sslmode=disable", hostPort.Port())

	dbpool, err := pgxpool.New(context.Background(), postgresUrl)
	if err != nil {
		t.Fatal(err)
		return nil, nil
	}

	byteArray, err := os.ReadFile("../schemas/postgres_v1.sql")
	if err != nil {
		t.Fatal(err)
		return nil, nil
	}

	_, err = dbpool.Exec(context.Background(), string(byteArray))
	if err != nil {
		t.Fatal(err)
		return nil, nil
	}

	return postgres, repository.TodoRepositoryImpl{DBPool: *dbpool}
}
*/

func SetupPostgres(t *testing.T) (tc.Container, common.TodoRepository) {
	dbname, user, password := "testdb", "user", "password"
	postgresPort := nat.Port("5432/tcp")
	postgres, err := tc.GenericContainer(context.Background(),
		tc.GenericContainerRequest{
			ContainerRequest: tc.ContainerRequest{
				Image:        "postgres:14.5",
				ExposedPorts: []string{postgresPort.Port()},
				Env: map[string]string{
					"POSTGRES_PASSWORD": password,
					"POSTGRES_USER":     user,
					"POSTGRES_DB":       dbname,
				},
				WaitingFor: wait.ForAll(
					wait.ForLog("database system is ready to accept connections"),
					wait.ForListeningPort(postgresPort),
				),
			},
			Started: true,
		})
	if err != nil {
		t.Fatal(err)
		return nil, nil
	}

	hostPort, err := postgres.MappedPort(context.Background(), postgresPort)
	if err != nil {
		t.Fatal(err)
		return nil, nil
	}

	postgresDataSourceName := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		"localhost", hostPort.Port(), user, password, dbname)

	dbpool, err := sql.Open("pgx", postgresDataSourceName)
	if err != nil {
		t.Fatal(err)
		return nil, nil
	}

	byteArray, err := os.ReadFile("../schemas/postgres_v1.sql")
	if err != nil {
		t.Fatal(err)
		return nil, nil
	}

	_, err = dbpool.Exec(string(byteArray))
	if err != nil {
		t.Fatal(err)
		return nil, nil
	}

	todoRepository, err := GetTodoRepository(dbpool)
	if err != nil {
		t.Fatal(err)
		return nil, nil
	}

	return postgres, todoRepository
}