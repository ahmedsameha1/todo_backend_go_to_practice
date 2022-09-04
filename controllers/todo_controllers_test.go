package controllers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/model"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var anErrString string = "An error"
var anError error = errors.New(anErrString)

func TestCreate(t *testing.T) {
	t.Run("Good case", func(t *testing.T) {
		gin.SetMode(gin.TestMode)
		mockCtrl := gomock.NewController(t)
		ginContextMock := common.NewMockWebContext(mockCtrl)
		todoRepositoryMock := common.NewMockTodoRepository(mockCtrl)
		errorHandlerMock := common.NewMockErrorHandler(mockCtrl)
		done := false
		todo := model.Todo{Title: "title1",
			Description: "description1",
			Done:        &done}
		ginContextMock.EXPECT().ShouldBindJSON(gomock.Any()).SetArg(0, todo)
		ginContextMock.EXPECT().JSON(http.StatusOK, todo)
		todoRepositoryMock.EXPECT().Create(&todo).Return(nil)
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
		todo := model.Todo{Title: "title1",
			Description: "description1",
			Done:        &done}
		ginContextMock.EXPECT().ShouldBindJSON(gomock.Any()).SetArg(0, todo)
		todoRepositoryMock.EXPECT().Create(&todo).Return(anError)
		errorHandlerMock.EXPECT().HandleAppError(anError, "", http.StatusInternalServerError)
		createTodo := Create(todoRepositoryMock, errorHandlerMock)
		assert.NotNil(t, createTodo)
		createTodo(ginContextMock)
	})

	t.Run("When required fields are not present in the web request body", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		ginContextMock := common.NewMockWebContext(mockCtrl)
		todoRepositoryMock := common.NewMockTodoRepository(mockCtrl)
		errorHandlerMock := common.NewMockErrorHandler(mockCtrl)
		ginContextMock.EXPECT().ShouldBindJSON(gomock.Any()).Return(anError)
		todoRepositoryMock.EXPECT().Create(gomock.Any()).Times(0)
		errorHandlerMock.EXPECT().HandleAppError(anError, "", http.StatusBadRequest)
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
		todos := []model.Todo{{Title: "title1", Description: "description1", Done: &todo1done},
			{Title: "title2", Description: "description2", Done: &todo2done},
			{Title: "title3", Description: "description3", Done: &todo3done}}
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
			Return(nil, anError)
		errorHandlerMock.EXPECT().HandleAppError(anError, "", http.StatusInternalServerError)
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
		todo := model.Todo{Title: "title1", Description: "description1", Done: &done, Id: todoId}
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
			return uuid.Nil, anError
		}
		todoRepositoryMock.EXPECT().GetById(gomock.Any()).Times(0)
		ginContextMock.EXPECT().Param("id")
		errorHandlerMock.EXPECT().HandleAppError(anError, "", http.StatusBadRequest)
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
		todoRepositoryMock.EXPECT().GetById(gomock.Any()).Return(nil, anError)
		ginContextMock.EXPECT().Param("id")
		errorHandlerMock.EXPECT().HandleAppError(anError, "", http.StatusInternalServerError)
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
		todos := []model.Todo{{Title: "title1", Description: "description1", Done: &todo1done},
			{Title: "title2", Description: "description2", Done: &todo2done},
			{Title: "title3", Description: "description3", Done: &todo3done}}
		todoRepositoryMock.EXPECT().GetAllByUserId(userId).Return(todos, nil)
		ginContextMock.EXPECT().Param("id").Return(userId.String())
		ginContextMock.EXPECT().JSON(http.StatusOK, todos)
		getAllByUserId(ginContextMock)
	})

	t.Run("When invalid user id is sent as a path parameter in the url", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		uUidParseMock := func(id string) (uuid.UUID, error) {
			return uuid.Nil, anError
		}
		todoRepositoryMock.EXPECT().GetAllByUserId(gomock.Any()).Times(0)
		getAllByUserId := GetAllByUserId(todoRepositoryMock, errorHandlerMock, uUidParseMock)
		ginContextMock.EXPECT().Param("id")
		errorHandlerMock.EXPECT().HandleAppError(anError, "", http.StatusBadRequest)
		assert.NotNil(t, getAllByUserId)
		getAllByUserId(ginContextMock)
	})

	t.Run("When TodoRepository returns an error", func(t *testing.T) {
		todoRepositoryMock, ginContextMock, errorHandlerMock := createMocks(t)
		userId := uuid.New()
		uUidParseMock := func(id string) (uuid.UUID, error) {
			return userId, nil
		}
		todoRepositoryMock.EXPECT().GetAllByUserId(gomock.Any()).Return(nil, anError)
		errorHandlerMock.EXPECT().HandleAppError(anError, "", http.StatusInternalServerError)
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
	gin.SetMode(gin.TestMode)
	todoRepositoryMock := createTodoRepositoryMock(t)
	update := Update(todoRepositoryMock)
	assert.NotNil(t, update)
	userId := uuid.New()
	done := false
	todo := model.Todo{Title: "title1",
		Description: "description1",
		Done:        &done}
	todoRepositoryMock.EXPECT().Update(&todo).Return(nil)
	todoJson, _ := json.Marshal(todo)
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	req, _ := http.NewRequest("PUT", "/todos/", bytes.NewBuffer(todoJson))
	ctx.Request = req
	ctx.Params = []gin.Param{{Key: "id", Value: userId.String()}}
	update(ctx)
	assert.Equal(t, http.StatusNoContent, r.Code)
	assert.Len(t, r.Body.Bytes(), 0)
}

func TestUpdateWhenInvalidIdIsSent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	todoRepositoryMock := createTodoRepositoryMock(t)
	update := Update(todoRepositoryMock)
	assert.NotNil(t, update)
	done := false
	todo := model.Todo{Title: "title1",
		Description: "description1",
		Done:        &done}
	todoRepositoryMock.EXPECT().Update(gomock.Any()).Times(0)
	todoJson, _ := json.Marshal(todo)
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	req, _ := http.NewRequest("PUT", "/todos/", bytes.NewBuffer(todoJson))
	ctx.Request = req
	ctx.Params = []gin.Param{{Key: "id", Value: "71ca04c4-2d88-4bc0-a5a3-47446098905n"}}
	update(ctx)
	assert.Equal(t, http.StatusBadRequest, r.Code)
	var got gin.H
	err := json.Unmarshal(r.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, got["error"], "UUID")
}

func TestUpdateRequiredFieldsAreNotPresentInRequestBody(t *testing.T) {
	gin.SetMode(gin.TestMode)
	todoRepositoryMock := createTodoRepositoryMock(t)
	update := Update(todoRepositoryMock)
	assert.NotNil(t, update)
	todoRepositoryMock.EXPECT().Update(gomock.Any()).Times(0)
	todoJson, _ := json.Marshal(model.Todo{})
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	req, _ := http.NewRequest("PUT", "/todos/", bytes.NewBuffer(todoJson))
	ctx.Request = req
	ctx.Params = []gin.Param{{Key: "id", Value: "71ca04c4-2d88-4bc0-a5a3-474460989058"}}
	update(ctx)
	assert.Equal(t, http.StatusBadRequest, r.Code)
	var got gin.H
	err := json.Unmarshal(r.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range []string{"Title", "Description", "Done"} {
		assert.Contains(t, got["error"], value)
	}
}

func TestUpdateWhenTodoRepositoryReturnAnError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	todoRepositoryMock := createTodoRepositoryMock(t)
	update := Update(todoRepositoryMock)
	assert.NotNil(t, update)
	userId := uuid.New()
	done := false
	todo := model.Todo{Title: "title1",
		Description: "description1",
		Done:        &done}
	todoRepositoryMock.EXPECT().Update(&todo).Return(errors.New("An error"))
	todoJson, _ := json.Marshal(todo)
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	req, _ := http.NewRequest("PUT", "/todos/", bytes.NewBuffer(todoJson))
	ctx.Request = req
	ctx.Params = []gin.Param{{Key: "id", Value: userId.String()}}
	update(ctx)
	assert.Equal(t, http.StatusInternalServerError, r.Code)
	var got gin.H
	err := json.Unmarshal(r.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, got["error"], "An error")
}

func TestDelete(t *testing.T) {
	gin.SetMode(gin.TestMode)
	todoRepositoryMock := createTodoRepositoryMock(t)
	todoId := uuid.New()
	todoRepositoryMock.EXPECT().Delete(todoId).Return(nil)
	delete := Delete(todoRepositoryMock)
	assert.NotNil(t, delete)
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	req, _ := http.NewRequest("DELETE", "/todos/", nil)
	ctx.Request = req
	ctx.Params = []gin.Param{{Key: "id", Value: todoId.String()}}
	delete(ctx)
	assert.Equal(t, http.StatusNoContent, r.Code)
	assert.Len(t, r.Body.Bytes(), 0)
}

func TestDeleteWhenInvalidIdIsSent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	todoRepositoryMock := createTodoRepositoryMock(t)
	todoRepositoryMock.EXPECT().Delete(gomock.Any()).Times(0)
	delete := Delete(todoRepositoryMock)
	assert.NotNil(t, delete)
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	req, _ := http.NewRequest("DELETE", "/todos/", nil)
	ctx.Request = req
	ctx.Params = []gin.Param{{Key: "id", Value: "71ca04c4-2d88-4bc0-a5a3-47446098905n"}}
	delete(ctx)
	assert.Equal(t, http.StatusBadRequest, r.Code)
	var got gin.H
	err := json.Unmarshal(r.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, got["error"], "UUID")
}

func TestDeleteWhenTodoRepositoryReturnAnError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	todoRepositoryMock := createTodoRepositoryMock(t)
	todoId := uuid.New()
	todoRepositoryMock.EXPECT().Delete(todoId).Return(anError)
	delete := Delete(todoRepositoryMock)
	assert.NotNil(t, delete)
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	req, _ := http.NewRequest("DELETE", "/todos/", nil)
	ctx.Request = req
	ctx.Params = []gin.Param{{Key: "id", Value: todoId.String()}}
	delete(ctx)
	assert.Equal(t, http.StatusInternalServerError, r.Code)
	var got gin.H
	err := json.Unmarshal(r.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	assert.Contains(t, got["error"], anErrString)
}

func createTodoRepositoryMock(t *testing.T) *common.MockTodoRepository {
	t.Helper()
	mockCtrl := gomock.NewController(t)
	return common.NewMockTodoRepository(mockCtrl)
}

func createMocks(t *testing.T) (*common.MockTodoRepository, *common.MockWebContext, *common.MockErrorHandler) {
	t.Helper()
	mockCtrl := gomock.NewController(t)
	return common.NewMockTodoRepository(mockCtrl), common.NewMockWebContext(mockCtrl), common.NewMockErrorHandler(mockCtrl)
}

func ExpectsErrorsGetVerb(t *testing.T, handler gin.HandlerFunc, path string, params []gin.Param, code int, errorContained []string) {
	t.Helper()
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	req, _ := http.NewRequest("GET", path, nil)
	ctx.Request = req
	if params != nil {
		ctx.Params = params
	}
	handler(ctx)
	assert.Equal(t, code, r.Code)
	var got gin.H
	err := json.Unmarshal(r.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range errorContained {
		assert.Contains(t, got["error"], value)
	}
}
