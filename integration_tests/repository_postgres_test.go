package integration_tests

import (
	"context"
	"fmt"
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

func setupPostgres(t *testing.T) (tc.Container, common.TodoRepository) {
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

func TestTodoRepositoryImplOnPostgres1(t *testing.T) {
	t.Run("Test todo creation", func(t *testing.T) {
		container, todoRepository := setupPostgres(t)
		defer container.Terminate(context.Background())
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

func TestTodoRepositoryImplOnPostgres2(t *testing.T) {
	t.Run("Test todo update: user id is not the same", func(t *testing.T) {
		container, todoRepository := setupPostgres(t)
		defer container.Terminate(context.Background())
		todoDone1 := true
		ti1, _ := time.Parse(time.RFC3339, "2022-09-21T14:07:05.768Z")
		todoId := uuid.New().String()
		expectedTodo1 := model.Todo{
			Id:          todoId,
			Title:       "title1",
			Description: "description1",
			Done:        &todoDone1,
			CreatedAt:   ti1,
		}
		userId1 := uuid.New().String()
		err := todoRepository.Create(&expectedTodo1, userId1)
		assert.NoError(t, err)
		todoDone2 := false
		ti2, _ := time.Parse(time.RFC3339, "2021-09-21T14:07:05.768Z")
		expectedTodo2 := model.Todo{
			Id:          todoId,
			Title:       "title1updated",
			Description: "description1updated",
			Done:        &todoDone2,
			CreatedAt:   ti2,
		}
		userId2 := uuid.New().String()
		err = todoRepository.Update(&expectedTodo2, userId2)
		assert.NoError(t, err)
		returnedTodo, err := todoRepository.GetById(todoId, userId1)
		assert.NoError(t, err)
		assert.Equal(t, expectedTodo1.Id, returnedTodo.Id)
		assert.Equal(t, expectedTodo1.Title, returnedTodo.Title)
		assert.Equal(t, expectedTodo1.Description, returnedTodo.Description)
		assert.Equal(t, expectedTodo1.Done, returnedTodo.Done)
		assert.Equal(t, expectedTodo1.CreatedAt, returnedTodo.CreatedAt.UTC())
	})
}

func TestTodoRepositoryImplOnPostgres3(t *testing.T) {
	t.Run("Test todo update: todo id is not the same", func(t *testing.T) {
		container, todoRepository := setupPostgres(t)
		defer container.Terminate(context.Background())
		todoDone1 := true
		ti1, _ := time.Parse(time.RFC3339, "2022-09-21T14:07:05.768Z")
		todoId1 := uuid.New().String()
		expectedTodo1 := model.Todo{
			Id:          todoId1,
			Title:       "title1",
			Description: "description1",
			Done:        &todoDone1,
			CreatedAt:   ti1,
		}
		userId := uuid.New().String()
		err := todoRepository.Create(&expectedTodo1, userId)
		assert.NoError(t, err)
		todoDone2 := false
		ti2, _ := time.Parse(time.RFC3339, "2021-09-21T14:07:05.768Z")
		todoId2 := uuid.New().String()
		expectedTodo2 := model.Todo{
			Id:          todoId2,
			Title:       "title1updated",
			Description: "description1updated",
			Done:        &todoDone2,
			CreatedAt:   ti2,
		}
		err = todoRepository.Update(&expectedTodo2, userId)
		assert.NoError(t, err)
		returnedTodo, err := todoRepository.GetById(todoId1, userId)
		assert.NoError(t, err)
		assert.Equal(t, expectedTodo1.Id, returnedTodo.Id)
		assert.Equal(t, expectedTodo1.Title, returnedTodo.Title)
		assert.Equal(t, expectedTodo1.Description, returnedTodo.Description)
		assert.Equal(t, expectedTodo1.Done, returnedTodo.Done)
		assert.Equal(t, expectedTodo1.CreatedAt, returnedTodo.CreatedAt.UTC())
	})
}

func TestTodoRepositoryImplOnPostgres4(t *testing.T) {
	t.Run("Test todo update: Good case", func(t *testing.T) {
		container, todoRepository := setupPostgres(t)
		defer container.Terminate(context.Background())
		todoDone1 := true
		ti1, _ := time.Parse(time.RFC3339, "2022-09-21T14:07:05.768Z")
		todoId := uuid.New().String()
		expectedTodo1 := model.Todo{
			Id:          todoId,
			Title:       "title1",
			Description: "description1",
			Done:        &todoDone1,
			CreatedAt:   ti1,
		}
		userId := uuid.New().String()
		err := todoRepository.Create(&expectedTodo1, userId)
		assert.NoError(t, err)
		todoDone2 := false
		ti2, _ := time.Parse(time.RFC3339, "2021-09-21T14:07:05.768Z")
		expectedTodo2 := model.Todo{
			Id:          todoId,
			Title:       "title1updated",
			Description: "description1updated",
			Done:        &todoDone2,
			CreatedAt:   ti2,
		}
		err = todoRepository.Update(&expectedTodo2, userId)
		assert.NoError(t, err)
		returnedTodo, err := todoRepository.GetById(todoId, userId)
		assert.NoError(t, err)
		assert.Equal(t, expectedTodo2.Id, returnedTodo.Id)
		assert.Equal(t, expectedTodo2.Title, returnedTodo.Title)
		assert.Equal(t, expectedTodo2.Description, returnedTodo.Description)
		assert.Equal(t, expectedTodo2.Done, returnedTodo.Done)
		assert.Equal(t, expectedTodo2.CreatedAt, returnedTodo.CreatedAt.UTC())
	})
}


func TestTodoRepositoryImplOnPostgres5(t *testing.T) {
	t.Run("Test todo deletion: user id is not the same", func(t *testing.T) {
		container, todoRepository := setupPostgres(t)
		defer container.Terminate(context.Background())
		todoDone := true
		ti, _ := time.Parse(time.RFC3339, "2022-09-21T14:07:05.768Z")
		todoId := uuid.New().String()
		expectedTodo := model.Todo{
			Id:          todoId,
			Title:       "title1",
			Description: "description1",
			Done:        &todoDone,
			CreatedAt:   ti,
		}
		userId1 := uuid.New().String()
		userId2 := uuid.New().String()
		err := todoRepository.Create(&expectedTodo, userId1)
		assert.NoError(t, err)
		err = todoRepository.Delete(todoId, userId2)
		assert.NoError(t, err)
		returnedTodo, err := todoRepository.GetById(todoId, userId1)
		assert.NoError(t, err)
		assert.NotNil(t, returnedTodo)
	})
}

func TestTodoRepositoryImplOnPostgres6(t *testing.T) {
	t.Run("Test todo deletion: todo id is not the same", func(t *testing.T) {
		container, todoRepository := setupPostgres(t)
		defer container.Terminate(context.Background())
		todoDone := true
		ti, _ := time.Parse(time.RFC3339, "2022-09-21T14:07:05.768Z")
		todoId1 := uuid.New().String()
		expectedTodo := model.Todo{
			Id:          todoId1,
			Title:       "title1",
			Description: "description1",
			Done:        &todoDone,
			CreatedAt:   ti,
		}
		userId := uuid.New().String()
		todoId2 := uuid.New().String()
		err := todoRepository.Create(&expectedTodo, userId)
		assert.NoError(t, err)
		err = todoRepository.Delete(todoId2, userId)
		assert.NoError(t, err)
		returnedTodo, err := todoRepository.GetById(todoId1, userId)
		assert.NoError(t, err)
		assert.NotNil(t, returnedTodo)
	})
}

func TestTodoRepositoryImplOnPostgres7(t *testing.T) {
	t.Run("Test todo deletion: Good case", func(t *testing.T) {
		container, todoRepository := setupPostgres(t)
		defer container.Terminate(context.Background())
		todoDone := true
		ti, _ := time.Parse(time.RFC3339, "2022-09-21T14:07:05.768Z")
		todoId := uuid.New().String()
		expectedTodo := model.Todo{
			Id:          todoId,
			Title:       "title1",
			Description: "description1",
			Done:        &todoDone,
			CreatedAt:   ti,
		}
		userId := uuid.New().String()
		err := todoRepository.Create(&expectedTodo, userId)
		assert.NoError(t, err)
		err = todoRepository.Delete(todoId, userId)
		assert.NoError(t, err)
		returnedTodo, err := todoRepository.GetById(todoId, userId)
		assert.Equal(t, repository.ErrNotFound, err)
		assert.Nil(t, returnedTodo)
	})
}