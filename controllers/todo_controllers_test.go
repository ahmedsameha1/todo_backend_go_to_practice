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

func TestGetAll(t *testing.T) {
	t.Run("Good case: there are no todos", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		ginContextMock := common.NewMockWebContext(mockCtrl)
		todoRepositoryMock := common.NewMockTodoRepository(mockCtrl)
		errorHandlerMock := common.NewMockErrorHandler(mockCtrl)
		todoRepositoryMock.EXPECT().GetAll().Return([]model.Todo{}, nil)
		ginContextMock.EXPECT().JSON(http.StatusOK, []model.Todo{})
		getAll := GetAll(todoRepositoryMock, errorHandlerMock)
		assert.NotNil(t, getAll)
		getAll(ginContextMock)
	})

	t.Run("Good case: there is at least one todo", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		ginContextMock := common.NewMockWebContext(mockCtrl)
		todoRepositoryMock := common.NewMockTodoRepository(mockCtrl)
		errorHandlerMock := common.NewMockErrorHandler(mockCtrl)
		todo1done := false
		todo2done := true
		todo3done := false
		todos := []model.Todo{{Id: uuid.New().String(), Title: "title1", Description: "description1", Done: &todo1done},
			{Id: uuid.New().String(), Title: "title2", Description: "description2", Done: &todo2done},
			{Id: uuid.New().String(), Title: "title3", Description: "description3", Done: &todo3done}}
		ginContextMock.EXPECT().JSON(http.StatusOK, todos)
		todoRepositoryMock.EXPECT().GetAll().
			Return(todos, nil)
		getTodos := GetAll(todoRepositoryMock, errorHandlerMock)
		assert.NotNil(t, getTodos)
		getTodos(ginContextMock)
	})

	t.Run("When TodoRepository returns an error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		ginContextMock := common.NewMockWebContext(mockCtrl)
		todoRepositoryMock := common.NewMockTodoRepository(mockCtrl)
		errorHandlerMock := common.NewMockErrorHandler(mockCtrl)
		todoRepositoryMock.EXPECT().GetAll().
			Return(nil, common.ErrError)
		errorHandlerMock.EXPECT().HandleAppError(common.ErrError, "", http.StatusInternalServerError)
		getAll := GetAll(todoRepositoryMock, errorHandlerMock)
		assert.NotNil(t, getAll)
		getAll(ginContextMock)
	})
}

func TestGetById(t *testing.T) {
	t.Run("Good case", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		todoId := uuid.New()
		uUidParseMock := func(id string) (uuid.UUID, error) {
			return todoId, nil
		}
		done := false
		todo := model.Todo{Id: todoId.String(), Title: "title1", Description: "description1", Done: &done}
		todoRepositoryMock.EXPECT().GetById(todoId).Return(&todo, nil)
		ginContextMock.EXPECT().Param("id").Return(todoId.String())
		ginContextMock.EXPECT().JSON(http.StatusOK, &todo)
		getById := GetById(todoRepositoryMock, errorHandlerMock, uUidParseMock)
		assert.NotNil(t, getById)
		getById(ginContextMock)
	})

	t.Run("When invalid id is sent as a path parameter in the url", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		uUidParseMock := func(id string) (uuid.UUID, error) {
			return uuid.Nil, common.ErrError
		}
		todoRepositoryMock.EXPECT().GetById(gomock.Any()).Times(0)
		ginContextMock.EXPECT().Param("id")
		errorHandlerMock.EXPECT().HandleAppError(common.ErrError, "", http.StatusBadRequest)
		getById := GetById(todoRepositoryMock, errorHandlerMock, uUidParseMock)
		assert.NotNil(t, getById)
		getById(ginContextMock)
	})

	t.Run("When TodoRepository returns an Error", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		todoId := uuid.New()
		uUidParseMock := func(id string) (uuid.UUID, error) {
			return todoId, nil
		}
		todoRepositoryMock.EXPECT().GetById(gomock.Any()).Return(nil, common.ErrError)
		ginContextMock.EXPECT().Param("id")
		errorHandlerMock.EXPECT().HandleAppError(common.ErrError, "", http.StatusInternalServerError)
		getById := GetById(todoRepositoryMock, errorHandlerMock, uUidParseMock)
		assert.NotNil(t, getById)
		getById(ginContextMock)
	})

	t.Run("When parse is nil", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		errorHandlerMock.EXPECT().HandleAppError(ErrParseIsNil, "", http.StatusInternalServerError)
		getById := GetById(todoRepositoryMock, errorHandlerMock, nil)
		assert.NotNil(t, getById)
		getById(ginContextMock)
	})
}

func TestGetAllByUserId(t *testing.T) {
	t.Run("Good case", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		userId := uuid.New()
		uUidParseMock := func(id string) (uuid.UUID, error) {
			return userId, nil
		}
		getAllByUserId := GetAllByUserId(todoRepositoryMock, errorHandlerMock, uUidParseMock)
		assert.NotNil(t, getAllByUserId)
		todo1done := false
		todo2done := true
		todo3done := false
		todos := []model.Todo{{Id: uuid.New().String(), Title: "title1", Description: "description1", Done: &todo1done},
			{Id: uuid.New().String(), Title: "title2", Description: "description2", Done: &todo2done},
			{Id: uuid.New().String(), Title: "title3", Description: "description3", Done: &todo3done}}
		todoRepositoryMock.EXPECT().GetAllByUserId(userId).Return(todos, nil)
		ginContextMock.EXPECT().Param("id").Return(userId.String())
		ginContextMock.EXPECT().JSON(http.StatusOK, todos)
		getAllByUserId(ginContextMock)
	})

	t.Run("When invalid user id is sent as a path parameter in the url", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		uUidParseMock := func(id string) (uuid.UUID, error) {
			return uuid.Nil, common.ErrError
		}
		todoRepositoryMock.EXPECT().GetAllByUserId(gomock.Any()).Times(0)
		getAllByUserId := GetAllByUserId(todoRepositoryMock, errorHandlerMock, uUidParseMock)
		ginContextMock.EXPECT().Param("id")
		errorHandlerMock.EXPECT().HandleAppError(common.ErrError, "", http.StatusBadRequest)
		assert.NotNil(t, getAllByUserId)
		getAllByUserId(ginContextMock)
	})

	t.Run("When TodoRepository returns an error", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		userId := uuid.New()
		uUidParseMock := func(id string) (uuid.UUID, error) {
			return userId, nil
		}
		todoRepositoryMock.EXPECT().GetAllByUserId(gomock.Any()).Return(nil, common.ErrError)
		errorHandlerMock.EXPECT().HandleAppError(common.ErrError, "", http.StatusInternalServerError)
		ginContextMock.EXPECT().Param("id")
		getAllByUserId := GetAllByUserId(todoRepositoryMock, errorHandlerMock, uUidParseMock)
		assert.NotNil(t, getAllByUserId)
		getAllByUserId(ginContextMock)
	})

	t.Run("When parse is nil", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		errorHandlerMock.EXPECT().HandleAppError(ErrParseIsNil, "", http.StatusInternalServerError)
		getAllByUserId := GetAllByUserId(todoRepositoryMock, errorHandlerMock, nil)
		assert.NotNil(t, getAllByUserId)
		getAllByUserId(ginContextMock)
	})
}

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
