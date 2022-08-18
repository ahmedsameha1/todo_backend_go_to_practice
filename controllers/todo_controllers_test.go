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
	gin.SetMode(gin.TestMode)
	mockCtrl := gomock.NewController(t)
	todoRepositoryMock := common.NewMockTodoRepository(mockCtrl)
	done := false
	todo := model.Todo{Title: "title1",
		Description: "description1",
		Done:        &done}
	todoRepositoryMock.EXPECT().Create(&todo).Return(nil)
	createTodo := Create(todoRepositoryMock)
	assert.NotNil(t, createTodo)
	v, _ := json.Marshal(todo)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	r, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(v))
	c.Request = r
	createTodo(c)
	assert.Equal(t, http.StatusOK, w.Code)
	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	var want gin.H
	err = json.Unmarshal(v, &want)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, want, got)
}

func TestCreateWhenTodoRepositoryReturnError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockCtrl := gomock.NewController(t)
	todoRepositoryMock := common.NewMockTodoRepository(mockCtrl)
	done := false
	todo := model.Todo{Title: "title1",
		Description: "description1",
		Done:        &done}
	todoRepositoryMock.EXPECT().Create(&todo).Return(errors.New("An error"))
	createTodo := Create(todoRepositoryMock)
	assert.NotNil(t, createTodo)
	ExpectsErrors(t, createTodo, todo,
		"POST", "/todos", http.StatusInternalServerError,
		[]string{"An error"})
}

func TestCreateRequiredFieldsAreNotPresentInRequestBody(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	todoRepositoryMock := common.NewMockTodoRepository(mockCtrl)
	todoRepositoryMock.EXPECT().Create(gomock.Any()).Times(0)
	createTodo := Create(todoRepositoryMock)
	assert.NotNil(t, createTodo)
	ExpectsErrors(t, createTodo, model.Todo{},
		"POST", "/todos", http.StatusBadRequest, []string{"Title", "Description", "Done"})
}

func TestGetAllWhenTodoRepositoryReturnEmptyArray(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockCtrl := gomock.NewController(t)
	todoRepositoryMock := common.NewMockTodoRepository(mockCtrl)
	todoRepositoryMock.EXPECT().GetAll().Return([]model.Todo{}, nil)
	getTodos := GetAll(todoRepositoryMock)
	assert.NotNil(t, getTodos)
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	req, _ := http.NewRequest("GET", "/todos", nil)
	ctx.Request = req
	getTodos(ctx)
	assert.Equal(t, http.StatusOK, r.Code)
	var got []gin.H
	err := json.Unmarshal(r.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
}

func TestGetAllWhenTodoRepositoryReturnNonEmptyArray(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockCtrl := gomock.NewController(t)
	todoRepositoryMock := common.NewMockTodoRepository(mockCtrl)
	todo1done := false
	todo2done := true
	todo3done := false
	todos := []model.Todo{{Title: "title1", Description: "description1", Done: &todo1done},
		{Title: "title2", Description: "description2", Done: &todo2done},
		{Title: "title3", Description: "description3", Done: &todo3done}}
	todoRepositoryMock.EXPECT().GetAll().
		Return(todos, nil)
	getTodos := GetAll(todoRepositoryMock)
	assert.NotNil(t, getTodos)
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	req, _ := http.NewRequest("GET", "/todos", nil)
	ctx.Request = req
	getTodos(ctx)
	assert.Equal(t, http.StatusOK, r.Code)
	var got []model.Todo
	err := json.Unmarshal(r.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	var want []model.Todo = todos
	assert.Equal(t, want, got)
}

func TestGetTodosWhenTodoRepositoryReturnError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockCtrl := gomock.NewController(t)
	todoRepositoryMock := common.NewMockTodoRepository(mockCtrl)
	todoRepositoryMock.EXPECT().GetAll().
		Return(nil, errors.New("An error"))
	getTodos := GetAll(todoRepositoryMock)
	assert.NotNil(t, getTodos)
	ExpectsErrorsGetVerb(t, getTodos, "/todos", nil, http.StatusInternalServerError,
		[]string{"An error"})
}

func TestGetById(t *testing.T) {
	gin.SetMode(gin.TestMode)
	todoRepositoryMock := createTodoRepositoryMock(t)
	done := false
	todoId := uuid.New()
	todo := model.Todo{Title: "title1", Description: "description1", Done: &done, Id: todoId}
	todoRepositoryMock.EXPECT().GetById(todoId).Return(&todo, nil)
	getTodoById := GetById(todoRepositoryMock)
	assert.NotNil(t, getTodoById)
	r := httptest.NewRecorder()
	urlf := "/todos/"
	req, _ := http.NewRequest("GET", urlf, nil)
	ctx, _ := gin.CreateTestContext(r)
	ctx.Request = req
	ctx.Params = []gin.Param{{Key: "id", Value: todoId.String()}}
	getTodoById(ctx)
	assert.Equal(t, http.StatusOK, r.Code)
	var got model.Todo
	err := json.Unmarshal(r.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, todo, got)
}

func TestGetByIdWhenInvalidIdIsSent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	todoRepositoryMock := createTodoRepositoryMock(t)
	todoRepositoryMock.EXPECT().GetById(gomock.Any()).Times(0)
	getTodoById := GetById(todoRepositoryMock)
	assert.NotNil(t, getTodoById)
	ExpectsErrorsGetVerb(t, getTodoById, "/todos/",
		[]gin.Param{{Key: "id", Value: "71ca04c4-2d88-4bc0-a5a3-47446098905n"}},
		http.StatusBadRequest, []string{"UUID"})
}

func TestGetByIdWhenTodoRepositoryReturnError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	todoRepositoryMock := createTodoRepositoryMock(t)
	todoRepositoryMock.EXPECT().GetById(gomock.Any()).Return(nil, errors.New("An error"))
	getTodoById := GetById(todoRepositoryMock)
	assert.NotNil(t, getTodoById)
	ExpectsErrorsGetVerb(t, getTodoById, "/todos/",
		[]gin.Param{{Key: "id", Value: "71ca04c4-2d88-4bc0-a5a3-474460989058"}},
		http.StatusInternalServerError, []string{})
}

func TestGetAllByUserId(t *testing.T) {
	gin.SetMode(gin.TestMode)
	todoRepositoryMock := createTodoRepositoryMock(t)
	getAllByUserId := GetAllByUserId(todoRepositoryMock)
	assert.NotNil(t, getAllByUserId)
	todo1done := false
	todo2done := true
	todo3done := false
	todos := []model.Todo{{Title: "title1", Description: "description1", Done: &todo1done},
		{Title: "title2", Description: "description2", Done: &todo2done},
		{Title: "title3", Description: "description3", Done: &todo3done}}
	userId := uuid.New()
	todoRepositoryMock.EXPECT().GetAllByUserId(userId).Return(todos, nil)
	r := httptest.NewRecorder()
	urlf := "/todos/users/"
	req, _ := http.NewRequest("GET", urlf, nil)
	ctx, _ := gin.CreateTestContext(r)
	ctx.Request = req
	ctx.Params = []gin.Param{{Key: "id", Value: userId.String()}}
	getAllByUserId(ctx)
	assert.Equal(t, http.StatusOK, r.Code)
	var got []model.Todo
	err := json.Unmarshal(r.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	var want []model.Todo = todos
	assert.Equal(t, want, got)
}

func TestGetAllByUserIdWhenInvalidIdIsSent(t *testing.T) {
	gin.SetMode(gin.TestMode)
	todoRepositoryMock := createTodoRepositoryMock(t)
	todoRepositoryMock.EXPECT().GetAllByUserId(gomock.Any()).Times(0)
	getAllByUserId := GetAllByUserId(todoRepositoryMock)
	assert.NotNil(t, getAllByUserId)
	ExpectsErrorsGetVerb(t, getAllByUserId, "/todos/users/",
		[]gin.Param{{Key: "id", Value: "71ca04c4-2d88-4bc0-a5a3-47446098905n"}},
		http.StatusBadRequest, []string{"UUID"})
}

func TestGetAllByUserIdWhenTodoRepositoryReturnError(t *testing.T) {
	gin.SetMode(gin.TestMode)
	todoRepositoryMock := createTodoRepositoryMock(t)
	todoRepositoryMock.EXPECT().GetAllByUserId(gomock.Any()).Return(nil, errors.New("An error"))
	getAllByUserId := GetAllByUserId(todoRepositoryMock)
	assert.NotNil(t, getAllByUserId)
	ExpectsErrorsGetVerb(t, getAllByUserId, "/todos/users/",
		[]gin.Param{{Key: "id", Value: "71ca04c4-2d88-4bc0-a5a3-474460989058"}},
		http.StatusInternalServerError, []string{})
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

func ExpectsErrors(t *testing.T, handler gin.HandlerFunc, td model.Todo,
	verb string, path string, code int, errorContained []string) {
	t.Helper()
	tdr, err := json.Marshal(td)
	if err != nil {
		t.Fatal(err)
	}
	r := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(r)
	req, _ := http.NewRequest(verb, path, bytes.NewBuffer(tdr))
	ctx.Request = req
	handler(ctx)
	assert.Equal(t, code, r.Code)
	var got gin.H
	err = json.Unmarshal(r.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range errorContained {
		assert.Contains(t, got["error"], value)
	}
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
