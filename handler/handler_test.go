package handler

import (
	"errors"
	"net/http"
	"testing"
	"time"

	"firebase.google.com/go/v4/auth"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/middleware"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/model"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestHandleError(t *testing.T) {
	t.Run("Normal case", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		mockCtrl := gomock.NewController(t)
		loggerMock := common.NewMockLogger(mockCtrl)
		ginContextMock := common.NewMockWebContext(mockCtrl)
		handlerErr := errors.New("handlerErr1")
		code := 404
		loggerMock.EXPECT().Printf("%v\n", handlerErr)
		ginContextMock.EXPECT().JSON(code, handlerErr)
		errorHandlerImpl := ErrorHandlerImpl{ Logger: loggerMock}
		errorHandlerImpl.HandleAppError(ginContextMock, handlerErr, code)
	})

	t.Run("When WebContext or Logger is nil, I trust that the app will panic!!", func(t *testing.T) {})
}

func TestCreate(t *testing.T) {
	t.Run("Good case", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		mockCtrl := gomock.NewController(t)
		ginContextMock := common.NewMockWebContext(mockCtrl)
		todoRepositoryMock := common.NewMockTodoRepository(mockCtrl)
		errorHandlerMock := common.NewMockErrorHandler(mockCtrl)
		done := false
		token := &auth.Token{UID: "sfweo"}
		todo := model.Todo{Id: uuid.New().String(), Title: "title1",
			Description: "description1",
			Done:        &done, CreatedAt: time.Now()}
		ginContextMock.EXPECT().ShouldBindJSON(gomock.Any()).SetArg(0, todo)
		ginContextMock.EXPECT().JSON(http.StatusOK, todo)
		ginContextMock.EXPECT().Get(middleware.AuthToken).Return(token, true)
		todoRepositoryMock.EXPECT().Create(&todo, token.UID).Return(nil)
		createTodo := Create(todoRepositoryMock, errorHandlerMock)
		assert.NotNil(t, createTodo)
		createTodo(ginContextMock)
	})

	t.Run("When TodoRepository returns an error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		ginContextMock := common.NewMockWebContext(mockCtrl)
		todoRepositoryMock := common.NewMockTodoRepository(mockCtrl)
		errorHandlerMock := common.NewMockErrorHandler(mockCtrl)
		done := false
		token := &auth.Token{UID: "sfweo"}
		todo := model.Todo{Id: uuid.New().String(), Title: "title1",
			Description: "description1",
			Done:        &done, CreatedAt: time.Now()}
		ginContextMock.EXPECT().ShouldBindJSON(gomock.Any()).SetArg(0, todo)
		ginContextMock.EXPECT().Get(middleware.AuthToken).Return(token, true)
		todoRepositoryMock.EXPECT().Create(&todo, token.UID).Return(common.ErrError)
		errorHandlerMock.EXPECT().HandleAppError(ginContextMock, common.ErrError, http.StatusInternalServerError)
		createTodo := Create(todoRepositoryMock, errorHandlerMock)
		assert.NotNil(t, createTodo)
		createTodo(ginContextMock)
	})

	t.Run("When required fields are not present in the web request body", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		ginContextMock := common.NewMockWebContext(mockCtrl)
		todoRepositoryMock := common.NewMockTodoRepository(mockCtrl)
		errorHandlerMock := common.NewMockErrorHandler(mockCtrl)
		token := &auth.Token{UID: "sfweo"}
		ginContextMock.EXPECT().ShouldBindJSON(gomock.Any()).Return(common.ErrError)
		ginContextMock.EXPECT().Get(middleware.AuthToken).Return(token, true)
		todoRepositoryMock.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
		errorHandlerMock.EXPECT().HandleAppError(ginContextMock, common.ErrError, http.StatusBadRequest)
		createTodo := Create(todoRepositoryMock, errorHandlerMock)
		assert.NotNil(t, createTodo)
		createTodo(ginContextMock)
	})

	t.Run("When there is no auth token in the web context", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		ginContextMock := common.NewMockWebContext(mockCtrl)
		todoRepositoryMock := common.NewMockTodoRepository(mockCtrl)
		errorHandlerMock := common.NewMockErrorHandler(mockCtrl)
		ginContextMock.EXPECT().Get(middleware.AuthToken).Return(nil, false)
		errorHandlerMock.EXPECT().HandleAppError(ginContextMock, middleware.ErrNoUID, http.StatusUnauthorized)
		createTodo := Create(todoRepositoryMock, errorHandlerMock)
		assert.NotNil(t, createTodo)
		createTodo(ginContextMock)
	})
}

func TestGetAll(t *testing.T) {
	t.Run("Good case: there are no todos", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		ginContextMock := common.NewMockWebContext(mockCtrl)
		todoRepositoryMock := common.NewMockTodoRepository(mockCtrl)
		errorHandlerMock := common.NewMockErrorHandler(mockCtrl)
		token := &auth.Token{UID: "wbfewh"}
		ginContextMock.EXPECT().Get(middleware.AuthToken).Return(token, true)
		todoRepositoryMock.EXPECT().GetAll(token.UID).Return([]model.Todo{}, nil)
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
		token := &auth.Token{UID: "wbfewh"}
		todo1done := false
		todo2done := true
		todo3done := false
		todos := []model.Todo{{Id: uuid.New().String(), Title: "title1", Description: "description1", Done: &todo1done, CreatedAt: time.Now()},
			{Id: uuid.New().String(), Title: "title2", Description: "description2", Done: &todo2done, CreatedAt: time.Now()},
			{Id: uuid.New().String(), Title: "title3", Description: "description3", Done: &todo3done, CreatedAt: time.Now()}}
		ginContextMock.EXPECT().Get(middleware.AuthToken).Return(token, true)
		ginContextMock.EXPECT().JSON(http.StatusOK, todos)
		todoRepositoryMock.EXPECT().GetAll(token.UID).Return(todos, nil)
		getTodos := GetAll(todoRepositoryMock, errorHandlerMock)
		assert.NotNil(t, getTodos)
		getTodos(ginContextMock)
	})

	t.Run("When TodoRepository returns an error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		ginContextMock := common.NewMockWebContext(mockCtrl)
		todoRepositoryMock := common.NewMockTodoRepository(mockCtrl)
		errorHandlerMock := common.NewMockErrorHandler(mockCtrl)
		token := &auth.Token{UID: "wbfewh"}
		ginContextMock.EXPECT().Get(middleware.AuthToken).Return(token, true)
		todoRepositoryMock.EXPECT().GetAll(token.UID).Return(nil, common.ErrError)
		errorHandlerMock.EXPECT().HandleAppError(ginContextMock, common.ErrError, http.StatusInternalServerError)
		getAll := GetAll(todoRepositoryMock, errorHandlerMock)
		assert.NotNil(t, getAll)
		getAll(ginContextMock)
	})

	t.Run("When there is no auth token in the web context", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		ginContextMock := common.NewMockWebContext(mockCtrl)
		todoRepositoryMock := common.NewMockTodoRepository(mockCtrl)
		errorHandlerMock := common.NewMockErrorHandler(mockCtrl)
		ginContextMock.EXPECT().Get(middleware.AuthToken).Return(nil, false)
		errorHandlerMock.EXPECT().HandleAppError(ginContextMock, middleware.ErrNoUID, http.StatusUnauthorized)
		getAll := GetAll(todoRepositoryMock, errorHandlerMock)
		assert.NotNil(t, getAll)
		getAll(ginContextMock)
	})
}

func TestGetById(t *testing.T) {
	t.Run("Good case", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		todoId := uuid.New()
		token := &auth.Token{UID: "heowh"}
		uUidParseMock := func(id string) (uuid.UUID, error) {
			return todoId, nil
		}
		done := false
		todo := model.Todo{Id: todoId.String(), Title: "title1", Description: "description1", Done: &done, CreatedAt: time.Now()}
		todoRepositoryMock.EXPECT().GetById(todoId.String(), token.UID).Return(&todo, nil)
		ginContextMock.EXPECT().Get(middleware.AuthToken).Return(token, true)
		ginContextMock.EXPECT().Param("id").Return(todoId.String())
		ginContextMock.EXPECT().JSON(http.StatusOK, &todo)
		getById := GetById(todoRepositoryMock, errorHandlerMock, uUidParseMock)
		assert.NotNil(t, getById)
		getById(ginContextMock)
	})

	t.Run("When invalid id is sent as a path parameter in the url", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		token := &auth.Token{UID: "heowh"}
		uUidParseMock := func(id string) (uuid.UUID, error) {
			return uuid.Nil, common.ErrError
		}
		todoRepositoryMock.EXPECT().GetById(gomock.Any(), token.UID).Times(0)
		ginContextMock.EXPECT().Param("id")
		ginContextMock.EXPECT().Get(middleware.AuthToken).Return(token, true)
		errorHandlerMock.EXPECT().HandleAppError(ginContextMock, common.ErrError, http.StatusBadRequest)
		getById := GetById(todoRepositoryMock, errorHandlerMock, uUidParseMock)
		assert.NotNil(t, getById)
		getById(ginContextMock)
	})

	t.Run("When TodoRepository returns an Error", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		todoId := uuid.New()
		token := &auth.Token{UID: "heowh"}
		uUidParseMock := func(id string) (uuid.UUID, error) {
			return todoId, nil
		}
		todoRepositoryMock.EXPECT().GetById(gomock.Any(), token.UID).Return(nil, common.ErrError)
		ginContextMock.EXPECT().Get(middleware.AuthToken).Return(token, true)
		ginContextMock.EXPECT().Param("id")
		errorHandlerMock.EXPECT().HandleAppError(ginContextMock, common.ErrError, http.StatusInternalServerError)
		getById := GetById(todoRepositoryMock, errorHandlerMock, uUidParseMock)
		assert.NotNil(t, getById)
		getById(ginContextMock)
	})

	t.Run("When parse is nil", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		token := &auth.Token{UID: "heowh"}
		ginContextMock.EXPECT().Get(middleware.AuthToken).Return(token, true)
		errorHandlerMock.EXPECT().HandleAppError(ginContextMock, ErrParseIsNil, http.StatusInternalServerError)
		getById := GetById(todoRepositoryMock, errorHandlerMock, nil)
		assert.NotNil(t, getById)
		getById(ginContextMock)
	})

	t.Run("When there is no auth token in the web context", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		ginContextMock.EXPECT().Get(middleware.AuthToken).Return(nil, false)
		errorHandlerMock.EXPECT().HandleAppError(ginContextMock, middleware.ErrNoUID, http.StatusUnauthorized)
		getById := GetById(todoRepositoryMock, errorHandlerMock, nil)
		assert.NotNil(t, getById)
		getById(ginContextMock)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Good case", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		token := &auth.Token{UID: "nfwseo"}
		update := Update(todoRepositoryMock, errorHandlerMock)
		assert.NotNil(t, update)
		done := false
		todo := model.Todo{Id: uuid.New().String(), Title: "title1",
			Description: "description1",
			Done:        &done, CreatedAt: time.Now()}
		todoRepositoryMock.EXPECT().Update(&todo, token.UID).Return(nil)
		ginContextMock.EXPECT().Get(middleware.AuthToken).Return(token, true)
		ginContextMock.EXPECT().ShouldBindJSON(gomock.Any()).SetArg(0, todo)
		ginContextMock.EXPECT().JSON(http.StatusNoContent, map[string]any{})
		update(ginContextMock)
	})

	t.Run("When required fields are not present in the web request body", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		token := &auth.Token{UID: "nfwseo"}
		update := Update(todoRepositoryMock, errorHandlerMock)
		assert.NotNil(t, update)
		todoRepositoryMock.EXPECT().Update(gomock.Any(), token.UID).Times(0)
		ginContextMock.EXPECT().Get(middleware.AuthToken).Return(token, true)
		ginContextMock.EXPECT().ShouldBindJSON(gomock.Any()).Return(common.ErrError)
		errorHandlerMock.EXPECT().HandleAppError(ginContextMock, common.ErrError, http.StatusBadRequest)
		update(ginContextMock)
	})

	t.Run("When TodoRepository returns an error", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		token := &auth.Token{UID: "nfwseo"}
		update := Update(todoRepositoryMock, errorHandlerMock)
		assert.NotNil(t, update)
		done := false
		todo := model.Todo{Id: uuid.New().String(), Title: "title1",
			Description: "description1",
			Done:        &done, CreatedAt: time.Now()}
		todoRepositoryMock.EXPECT().Update(&todo, token.UID).Return(common.ErrError)
		ginContextMock.EXPECT().Get(middleware.AuthToken).Return(token, true)
		ginContextMock.EXPECT().ShouldBindJSON(gomock.Any()).SetArg(0, todo)
		errorHandlerMock.EXPECT().HandleAppError(ginContextMock, common.ErrError, http.StatusInternalServerError)
		update(ginContextMock)
	})

	t.Run("When there is no auth token in the web context", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		update := Update(todoRepositoryMock, errorHandlerMock)
		assert.NotNil(t, update)
		ginContextMock.EXPECT().Get(middleware.AuthToken).Return(nil, false)
		errorHandlerMock.EXPECT().HandleAppError(ginContextMock, middleware.ErrNoUID, http.StatusUnauthorized)
		update(ginContextMock)
	})
}

func TestDelete(t *testing.T) {
	t.Run("Good case", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		token := &auth.Token{UID: "oiwhbegfwh"}
		todoId := uuid.New()
		uUidParseMock := func(id string) (uuid.UUID, error) {
			return todoId, nil
		}
		todoRepositoryMock.EXPECT().Delete(todoId.String(), token.UID).Return(nil)
		ginContextMock.EXPECT().Get(middleware.AuthToken).Return(token, true)
		ginContextMock.EXPECT().Param("id").Return(todoId.String())
		ginContextMock.EXPECT().JSON(http.StatusNoContent, map[string]any{})
		delete := Delete(todoRepositoryMock, errorHandlerMock, uUidParseMock)
		assert.NotNil(t, delete)
		delete(ginContextMock)
	})

	t.Run("When invalid todo id is sent as a path parameter in the url", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		token := &auth.Token{UID: "oiwhbegfwh"}
		uUidParseMock := func(id string) (uuid.UUID, error) {
			return uuid.Nil, common.ErrError
		}
		todoRepositoryMock.EXPECT().Delete(gomock.Any(), token.UID).Times(0)
		ginContextMock.EXPECT().Get(middleware.AuthToken).Return(token, true)
		ginContextMock.EXPECT().Param("id")
		errorHandlerMock.EXPECT().HandleAppError(ginContextMock, common.ErrError, http.StatusBadRequest)
		delete := Delete(todoRepositoryMock, errorHandlerMock, uUidParseMock)
		assert.NotNil(t, delete)
		delete(ginContextMock)
	})

	t.Run("When TodoRepository return an error", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		token := &auth.Token{UID: "oiwhbegfwh"}
		todoId := uuid.New()
		uUidParseMock := func(id string) (uuid.UUID, error) {
			return todoId, nil
		}
		todoRepositoryMock.EXPECT().Delete(todoId.String(), token.UID).Return(common.ErrError)
		ginContextMock.EXPECT().Get(middleware.AuthToken).Return(token, true)
		ginContextMock.EXPECT().Param("id").Return(todoId.String())
		errorHandlerMock.EXPECT().HandleAppError(ginContextMock, common.ErrError, http.StatusInternalServerError)
		delete := Delete(todoRepositoryMock, errorHandlerMock, uUidParseMock)
		assert.NotNil(t, delete)
		delete(ginContextMock)
	})

	t.Run("When parse is nil", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		token := &auth.Token{UID: "oiwhbegfwh"}
		ginContextMock.EXPECT().Get(middleware.AuthToken).Return(token, true)
		errorHandlerMock.EXPECT().HandleAppError(ginContextMock, ErrParseIsNil, http.StatusInternalServerError)
		delete := Delete(todoRepositoryMock, errorHandlerMock, nil)
		delete(ginContextMock)
	})

	t.Run("When there is no auth token in the web context", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		ginContextMock.EXPECT().Get(middleware.AuthToken).Return(nil, false)
		errorHandlerMock.EXPECT().HandleAppError(ginContextMock, middleware.ErrNoUID, http.StatusUnauthorized)
		delete := Delete(todoRepositoryMock, errorHandlerMock, nil)
		delete(ginContextMock)
	})
}

func createMocks(t *testing.T) (*common.MockTodoRepository, *common.MockWebContext, *common.MockErrorHandler) {
	t.Helper()
	mockCtrl := gomock.NewController(t)
	return common.NewMockTodoRepository(mockCtrl), common.NewMockWebContext(mockCtrl), common.NewMockErrorHandler(mockCtrl)
}
