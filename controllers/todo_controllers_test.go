package controllers

import (
	"net/http"
	"testing"

	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/model"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUpdate(t *testing.T) {
	t.Run("Good case", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		update := Update(todoRepositoryMock, errorHandlerMock)
		assert.NotNil(t, update)
		done := false
		todo := model.Todo{Id: uuid.New().String(), Title: "title1",
			Description: "description1",
			Done:        &done}
		todoRepositoryMock.EXPECT().Update(&todo).Return(nil)
		ginContextMock.EXPECT().ShouldBindJSON(gomock.Any()).SetArg(0, todo)
		ginContextMock.EXPECT().JSON(http.StatusNoContent, map[string]any{})
		update(ginContextMock)
	})

	t.Run("When required fields are not present in the web request body", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		update := Update(todoRepositoryMock, errorHandlerMock)
		assert.NotNil(t, update)
		todoRepositoryMock.EXPECT().Update(gomock.Any()).Times(0)
		ginContextMock.EXPECT().ShouldBindJSON(gomock.Any()).Return(common.ErrError)
		errorHandlerMock.EXPECT().HandleAppError(common.ErrError, "", http.StatusBadRequest)
		update(ginContextMock)
	})

	t.Run("When TodoRepository returns an error", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		update := Update(todoRepositoryMock, errorHandlerMock)
		assert.NotNil(t, update)
		done := false
		todo := model.Todo{Id: uuid.New().String(), Title: "title1",
			Description: "description1",
			Done:        &done}
		todoRepositoryMock.EXPECT().Update(&todo).Return(common.ErrError)
		ginContextMock.EXPECT().ShouldBindJSON(gomock.Any()).SetArg(0, todo)
		errorHandlerMock.EXPECT().HandleAppError(common.ErrError, "", http.StatusInternalServerError)
		update(ginContextMock)
	})
}

func TestDelete(t *testing.T) {
	t.Run("Good case", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		todoId := uuid.New()
		uUidParseMock := func(id string) (uuid.UUID, error) {
			return todoId, nil
		}
		todoRepositoryMock.EXPECT().Delete(todoId).Return(nil)
		ginContextMock.EXPECT().Param("id").Return(todoId.String())
		ginContextMock.EXPECT().JSON(http.StatusNoContent, map[string]any{})
		delete := Delete(todoRepositoryMock, errorHandlerMock, uUidParseMock)
		assert.NotNil(t, delete)
		delete(ginContextMock)
	})

	t.Run("When invalid todo id is sent as a path parameter in the url", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		uUidParseMock := func(id string) (uuid.UUID, error) {
			return uuid.Nil, common.ErrError
		}
		todoRepositoryMock.EXPECT().Delete(gomock.Any()).Times(0)
		ginContextMock.EXPECT().Param("id")
		errorHandlerMock.EXPECT().HandleAppError(common.ErrError, "", http.StatusBadRequest)
		delete := Delete(todoRepositoryMock, errorHandlerMock, uUidParseMock)
		assert.NotNil(t, delete)
		delete(ginContextMock)
	})

	t.Run("When TodoRepository return an error", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		todoId := uuid.New()
		uUidParseMock := func(id string) (uuid.UUID, error) {
			return todoId, nil
		}
		todoRepositoryMock.EXPECT().Delete(todoId).Return(common.ErrError)
		ginContextMock.EXPECT().Param("id").Return(todoId.String())
		errorHandlerMock.EXPECT().HandleAppError(common.ErrError, "", http.StatusInternalServerError)
		delete := Delete(todoRepositoryMock, errorHandlerMock, uUidParseMock)
		assert.NotNil(t, delete)
		delete(ginContextMock)
	})

	t.Run("When parse is nil", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		errorHandlerMock.EXPECT().HandleAppError(ErrParseIsNil, "", http.StatusInternalServerError)
		delete := Delete(todoRepositoryMock, errorHandlerMock, nil)
		delete(ginContextMock)
	})
}

func createMocks(t *testing.T) (*common.MockTodoRepository, *common.MockWebContext, *common.MockErrorHandler) {
	t.Helper()
	mockCtrl := gomock.NewController(t)
	return common.NewMockTodoRepository(mockCtrl), common.NewMockWebContext(mockCtrl), common.NewMockErrorHandler(mockCtrl)
}
