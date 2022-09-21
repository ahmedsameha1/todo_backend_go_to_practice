package integration_tests

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/model"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/repository"
	"github.com/docker/go-connections/nat"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	tc "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var todoRepository common.TodoRepository

func TestMain(m *testing.M) {
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
		log.Fatal("start:", err)
	}

	hostPort, err := postgres.MappedPort(context.Background(), postgresPort)
	if err != nil {
		log.Fatal("map:", err)
	}

	postgresUrl := fmt.Sprintf("postgres://user:password@localhost:%s?sslmode=disable", hostPort.Port())

	dbpool, err := pgxpool.New(context.Background(), postgresUrl)
	if err != nil {
		log.Fatal("Unable to connect to the database", err)
	}

	byteArray, err := os.ReadFile("../schemas/postgres_v1.sql")
	if err != nil {
		log.Fatal("read schema error", err)
	}

	_, err = dbpool.Exec(context.Background(), string(byteArray))
	if err != nil {
		log.Fatal("create schema error", err)
	}

	todoRepository = repository.TodoRepositoryImpl{DBPool: *dbpool}

	os.Exit(m.Run())

}

func TestTodoRepositoryImplOnPostgres(t *testing.T) {
	t.Run("Test todo creation", func(t *testing.T) {
		todoDone := false
		ti, _ := time.Parse(time.RFC3339, "2022-09-21T14:07:05.768Z")
		expectedTodo := model.Todo{
			Id:          uuid.New().String(),
			Title:       "title1",
			Description: "description1",
			Done:        &todoDone,
			CreatedAt:   ti,
		}
		userId := uuid.New().String()
		err := todoRepository.Create(&expectedTodo, userId)
		assert.NoError(t, err)
		todos, err := todoRepository.GetAll(userId)
		assert.NoError(t, err)
		assert.Equal(t, len(todos), 1)
		returnedTodo := todos[0]
		assert.Equal(t, expectedTodo.Id, returnedTodo.Id)
		assert.Equal(t, expectedTodo.Title, returnedTodo.Title)
		assert.Equal(t, expectedTodo.Description, returnedTodo.Description)
		assert.Equal(t, expectedTodo.Done, returnedTodo.Done)
		assert.Equal(t, expectedTodo.CreatedAt, returnedTodo.CreatedAt.UTC())
	})
}
