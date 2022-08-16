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

func TestCreate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockCtrl := gomock.NewController(t)
	todoRepositoryMock := common.NewMockTodoRepository(mockCtrl)
	done := false
	todo := model.Todo{Title: "title1",
		Description: "description1",
		Done:        &done}
	todoRepositoryMock.EXPECT().Create(&todo).Return(&todo, nil)
	createTodo := Create(todoRepositoryMock)
	assert.NotNil(t, createTodo)
	///
	v, _ := json.Marshal(todo)
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	r, _ := http.NewRequest("POST", "/todos", bytes.NewBuffer(v))
	c.Request = r
	createTodo(c)
	if w.Code != http.StatusOK {
		t.Errorf("Expects response code to be: %d, but it is: %d", http.StatusOK, w.Code)
	}
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
	//todoRepositoryMock.On("CreateTodo", &todo).Return(nil, errors.New("An error"))
	todoRepositoryMock.EXPECT().Create(&todo).Return(nil, errors.New("An error"))
	createTodo := Create(todoRepositoryMock)
	assert.NotNil(t, createTodo)
	///
	//todoRepositoryMock.AssertExpectations(t)
	ExpectsErrors(t, createTodo, todo,
		"POST", "/todos", http.StatusInternalServerError,
		[]string{"An error"})
}
func TestCreateNoTitleInRequestBody(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	todoRepositoryMock := common.NewMockTodoRepository(mockCtrl)
	createTodo := Create(todoRepositoryMock)
	assert.NotNil(t, createTodo)
	///
	ExpectsErrors(t, createTodo, model.Todo{},
		"POST", "/todos", http.StatusBadRequest, []string{"Title"})
}

func TestCreateEmptyTitleInRequestBody(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	todoRepositoryMock := common.NewMockTodoRepository(mockCtrl)
	createTodo := Create(todoRepositoryMock)
	assert.NotNil(t, createTodo)
	///
	ExpectsErrors(t, createTodo, model.Todo{Title: ""},
		"POST", "/todos", http.StatusBadRequest, []string{"Title"})
}

func TestCreateNoDescriptionInRequestBody(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	todoRepositoryMock := common.NewMockTodoRepository(mockCtrl)
	createTodo := Create(todoRepositoryMock)
	assert.NotNil(t, createTodo)
	///
	ExpectsErrors(t, createTodo, model.Todo{Title: "title1"},
		"POST", "/todos", http.StatusBadRequest, []string{"Description"})
}

func TestCreateEmptyDescriptionInRequestBody(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	todoRepositoryMock := common.NewMockTodoRepository(mockCtrl)
	createTodo := Create(todoRepositoryMock)
	assert.NotNil(t, createTodo)
	///
	ExpectsErrors(t, createTodo, model.Todo{Title: "", Description: ""},
		"POST", "/todos", http.StatusBadRequest, []string{"Title"})
}

func TestCreateEmptyTitleAndNoDescriptionInRequestBody(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	todoRepositoryMock := common.NewMockTodoRepository(mockCtrl)
	createTodo := Create(todoRepositoryMock)
	assert.NotNil(t, createTodo)
	///
	ExpectsErrors(t, createTodo, model.Todo{Title: ""},
		"POST", "/todos", http.StatusBadRequest, []string{"Title", "Description"})
}

func TestCreateNoDoneInRequestBody(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	todoRepositoryMock := common.NewMockTodoRepository(mockCtrl)
	createTodo := Create(todoRepositoryMock)
	assert.NotNil(t, createTodo)
	///
	ExpectsErrors(t, createTodo, model.Todo{Title: "title1",
		Description: "description1"},
		"POST", "/todos", http.StatusBadRequest, []string{"Done"})
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
	if r.Code != http.StatusOK {
		t.Errorf("Expects response code to be: %d, but it is: %d", http.StatusOK, r.Code)
	}
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
	if r.Code != http.StatusOK {
		t.Errorf("Expects response code to be: %d, but it is: %d", http.StatusOK, r.Code)
	}
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
	if r.Code != code {
		t.Errorf("Code should be: %d but it is: %d", code, r.Code)
	}
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
	if r.Code != code {
		t.Errorf("Code should be: %d but it is: %d", code, r.Code)
	}
	var got gin.H
	err := json.Unmarshal(r.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	for _, value := range errorContained {
		assert.Contains(t, got["error"], value)
	}
}
