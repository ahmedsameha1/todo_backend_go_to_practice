package integration_tests

import (
	"context"
	"testing"
	"time"

	"github.com/ahmedsameha1/todo_backend_go_to_practice/model"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/repository"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
)

func TestTodoRepositoryImplOnPostgres1(t *testing.T) {
	t.Run("Test todo creation", func(t *testing.T) {
		container, todoRepository := repository.SetupPostgres(t)
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
		assert.Equal(t, expectedTodo, returnedTodo)
	})
}

func TestTodoRepositoryImplOnPostgres2(t *testing.T) {
	t.Run("Test todo update: user id is not the same", func(t *testing.T) {
		container, todoRepository := repository.SetupPostgres(t)
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
		assert.Equal(t, &expectedTodo1, returnedTodo)
	})
}

func TestTodoRepositoryImplOnPostgres3(t *testing.T) {
	t.Run("Test todo update: todo id is not the same", func(t *testing.T) {
		container, todoRepository := repository.SetupPostgres(t)
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
		assert.Equal(t, &expectedTodo1, returnedTodo)
	})
}

func TestTodoRepositoryImplOnPostgres4(t *testing.T) {
	t.Run("Test todo update: Good case", func(t *testing.T) {
		container, todoRepository := repository.SetupPostgres(t)
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
		assert.Equal(t, &expectedTodo2, returnedTodo)
	})
}

func TestTodoRepositoryImplOnPostgres5(t *testing.T) {
	t.Run("Test todo deletion: user id is not the same", func(t *testing.T) {
		container, todoRepository := repository.SetupPostgres(t)
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
		container, todoRepository := repository.SetupPostgres(t)
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
		container, todoRepository := repository.SetupPostgres(t)
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

func TestTodoRepositoryImplOnPostgres8(t *testing.T) {
	t.Run("Test GetAll", func(t *testing.T) {
		container, todoRepository := repository.SetupPostgres(t)
		defer container.Terminate(context.Background())
		todoDone1 := false
		ti1, _ := time.Parse(time.RFC3339, "2020-09-21T14:07:05.768Z")
		todoId1 := uuid.New().String()
		expectedTodo1 := model.Todo{
			Id:          todoId1,
			Title:       "title1",
			Description: "description1",
			Done:        &todoDone1,
			CreatedAt:   ti1,
		}
		todoDone2 := false
		ti2, _ := time.Parse(time.RFC3339, "2021-09-21T14:07:05.768Z")
		todoId2 := uuid.New().String()
		expectedTodo2 := model.Todo{
			Id:          todoId2,
			Title:       "title2",
			Description: "description2",
			Done:        &todoDone2,
			CreatedAt:   ti2,
		}
		todoDone3 := true
		ti3, _ := time.Parse(time.RFC3339, "2022-09-21T14:07:05.768Z")
		todoId3 := uuid.New().String()
		expectedTodo3 := model.Todo{
			Id:          todoId3,
			Title:       "title3",
			Description: "description3",
			Done:        &todoDone3,
			CreatedAt:   ti3,
		}
		todoDone4 := false
		ti4, _ := time.Parse(time.RFC3339, "2022-07-21T14:07:05.768Z")
		todoId4 := uuid.New().String()
		expectedTodo4 := model.Todo{
			Id:          todoId4,
			Title:       "title4",
			Description: "description4",
			Done:        &todoDone4,
			CreatedAt:   ti4,
		}
		userId1 := uuid.New().String()
		userId2 := uuid.New().String()
		todoRepository.Create(&expectedTodo2, userId1)
		todoRepository.Create(&expectedTodo1, userId1)
		todoRepository.Create(&expectedTodo4, userId2)
		todoRepository.Create(&expectedTodo3, userId1)
		returnedTodos, err := todoRepository.GetAll(userId1)
		assert.NoError(t, err)
		assert.Equal(t, expectedTodo3, returnedTodos[0])
		assert.Equal(t, expectedTodo2, returnedTodos[1])
		assert.Equal(t, expectedTodo1, returnedTodos[2])
		returnedTodos, err = todoRepository.GetAll(userId2)
		assert.NoError(t, err)
		assert.Equal(t, expectedTodo4, returnedTodos[0])
	})
}
