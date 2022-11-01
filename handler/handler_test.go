package handler

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
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
		http_recorder := httptest.NewRecorder()
		gin_context, _ := gin.CreateTestContext(http_recorder)
		handlerErr := errors.New("handlerErr1")
		code := 404
		loggerMock.EXPECT().Printf("%v\n", handlerErr)
		errorHandlerImpl := ErrorHandlerImpl{Logger: loggerMock}
		errorHandlerImpl.HandleAppError(gin_context, handlerErr, code)
		assert.Equal(t, code, http_recorder.Code)
		var got gin.H
		err := json.Unmarshal(http_recorder.Body.Bytes(), &got)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, gin.H{"error": handlerErr.Error()}, got)
	})

	t.Run("When WebContext or Logger is nil, I trust that the app will panic!!", func(t *testing.T) {})
}

func TestCreate(t *testing.T) {
	t.Run("Good case", func(t *testing.T) {
		todoRepositoryMock, gin_context, http_recorder, errorHandlerMock := createMocks(t)
		createTodo := Create(todoRepositoryMock, errorHandlerMock)
		done := false
		token := &auth.Token{UID: "sfweo"}
		ti, _ := time.Parse(time.RFC3339, "2022-09-21T14:07:05.768Z")
		todo := model.Todo{Id: uuid.New().String(), Title: "title1",
			Description: "description1",
			Done:        &done, CreatedAt: ti}
		todoRepositoryMock.EXPECT().Create(&todo, token.UID).Return(nil)
		json_bytes, err := json.Marshal(todo)
		if err != nil {
			t.Fatal(err)
		}
		web_request := &http.Request{
			Body:   io.NopCloser(bytes.NewBuffer(json_bytes)),
			Header: map[string][]string{"Content-Type": {"application/json"}}}
		gin_context.Request = web_request
		gin_context.Set(middleware.AuthToken, token)
		createTodo(gin_context)
		assert.Equal(t, http.StatusOK, http_recorder.Code)
		var got model.Todo
		err = json.Unmarshal(http_recorder.Body.Bytes(), &got)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, todo, got)
	})

	t.Run("When TodoRepository returns an error", func(t *testing.T) {
		todoRepositoryMock, gin_context, _, errorHandlerMock := createMocks(t)
		createTodo := Create(todoRepositoryMock, errorHandlerMock)
		done := false
		token := &auth.Token{UID: "sfweo"}
		ti, _ := time.Parse(time.RFC3339, "2022-09-21T14:07:05.768Z")
		todo := model.Todo{Id: uuid.New().String(), Title: "title1",
			Description: "description1",
			Done:        &done, CreatedAt: ti}
		json_bytes, err := json.Marshal(todo)
		if err != nil {
			t.Fatal(err)
		}
		web_request := &http.Request{
			Body:   io.NopCloser(bytes.NewBuffer(json_bytes)),
			Header: map[string][]string{"Content-Type": {"application/json"}}}
		gin_context.Request = web_request
		gin_context.Set(middleware.AuthToken, token)
		todoRepositoryMock.EXPECT().Create(&todo, token.UID).Return(common.ErrError)
		errorHandlerMock.EXPECT().HandleAppError(gin_context, common.ErrError, http.StatusInternalServerError)
		createTodo(gin_context)
	})

	t.Run("When required fields are not present in the web request body", func(t *testing.T) {
		todoRepositoryMock, gin_context, _, errorHandlerMock := createMocks(t)
		createTodo := Create(todoRepositoryMock, errorHandlerMock)
		done := false
		token := &auth.Token{UID: "sfweo"}
		ti, _ := time.Parse(time.RFC3339, "2022-09-21T14:07:05.768Z")
		todo := model.Todo{Id: uuid.New().String(), Title: "title1",
			Done: &done, CreatedAt: ti}
		json_bytes, err := json.Marshal(todo)
		if err != nil {
			t.Fatal(err)
		}
		web_request := &http.Request{
			Body:   io.NopCloser(bytes.NewBuffer(json_bytes)),
			Header: map[string][]string{"Content-Type": {"application/json"}}}
		gin_context.Request = web_request
		gin_context.Set(middleware.AuthToken, token)
		todoRepositoryMock.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
		errorHandlerMock.EXPECT().HandleAppError(gin_context, gomock.Any(), http.StatusBadRequest).
			DoAndReturn(func(ctx *gin.Context, err error, code int) {
				if !strings.Contains(err.Error(), "Description") {
					t.Fail()
				}
			})
		createTodo(gin_context)
	})

	t.Run("When there is no auth token in the web context", func(t *testing.T) {
		todoRepositoryMock, gin_context, _, errorHandlerMock := createMocks(t)
		errorHandlerMock.EXPECT().HandleAppError(gin_context, middleware.ErrNoUID, http.StatusUnauthorized)
		createTodo := Create(todoRepositoryMock, errorHandlerMock)
		done := false
		ti, _ := time.Parse(time.RFC3339, "2022-09-21T14:07:05.768Z")
		todo := model.Todo{Id: uuid.New().String(), Title: "title1",
			Description: "description1",
			Done:        &done, CreatedAt: ti}
		todoRepositoryMock.EXPECT().Create(gomock.Any(), gomock.Any()).Times(0)
		json_bytes, err := json.Marshal(todo)
		if err != nil {
			t.Fatal(err)
		}
		web_request := &http.Request{
			Body:   io.NopCloser(bytes.NewBuffer(json_bytes)),
			Header: map[string][]string{"Content-Type": {"application/json"}}}
		gin_context.Request = web_request
		createTodo(gin_context)
	})
}

func TestGetAll(t *testing.T) {
	t.Run("Good case: there are no todos", func(t *testing.T) {
		todoRepositoryMock, gin_context, http_recorder, errorHandlerMock := createMocks(t)
		token := &auth.Token{UID: "wbfewh"}
		gin_context.Set(middleware.AuthToken, token)
		todoRepositoryMock.EXPECT().GetAll(token.UID).Return([]model.Todo{}, nil)
		getAll := GetAll(todoRepositoryMock, errorHandlerMock)
		getAll(gin_context)
		assert.Equal(t, http.StatusOK, http_recorder.Code)
		var got []model.Todo
		err := json.Unmarshal(http_recorder.Body.Bytes(), &got)
		if err != nil {
			t.Fatal(err)
		}
		assert.Empty(t, got)
	})

	t.Run("Good case: there is at least one todo", func(t *testing.T) {
		todoRepositoryMock, gin_context, http_recorder, errorHandlerMock := createMocks(t)
		token := &auth.Token{UID: "wbfewh"}
		todo1done := false
		todo2done := true
		todo3done := false
		ti1, _ := time.Parse(time.RFC3339, "2022-07-21T14:07:05.768Z")
		ti2, _ := time.Parse(time.RFC3339, "2022-08-21T14:07:05.768Z")
		ti3, _ := time.Parse(time.RFC3339, "2022-09-21T14:07:05.768Z")
		todos := []model.Todo{{Id: uuid.New().String(), Title: "title1", Description: "description1", Done: &todo1done, CreatedAt: ti1},
			{Id: uuid.New().String(), Title: "title2", Description: "description2", Done: &todo2done, CreatedAt: ti2},
			{Id: uuid.New().String(), Title: "title3", Description: "description3", Done: &todo3done, CreatedAt: ti3}}
		gin_context.Set(middleware.AuthToken, token)
		todoRepositoryMock.EXPECT().GetAll(token.UID).Return(todos, nil)
		getAll := GetAll(todoRepositoryMock, errorHandlerMock)
		getAll(gin_context)
		assert.Equal(t, http.StatusOK, http_recorder.Code)
		var got []model.Todo
		err := json.Unmarshal(http_recorder.Body.Bytes(), &got)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, todos, got)
	})

	t.Run("When TodoRepository returns an error", func(t *testing.T) {
		todoRepositoryMock, gin_context, _, errorHandlerMock := createMocks(t)
		token := &auth.Token{UID: "wbfewh"}
		gin_context.Set(middleware.AuthToken, token)
		todoRepositoryMock.EXPECT().GetAll(token.UID).Return(nil, common.ErrError)
		errorHandlerMock.EXPECT().HandleAppError(gin_context, common.ErrError, http.StatusInternalServerError)
		getAll := GetAll(todoRepositoryMock, errorHandlerMock)
		getAll(gin_context)
	})

	t.Run("When there is no auth token in the web context", func(t *testing.T) {
		todoRepositoryMock, gin_context, _, errorHandlerMock := createMocks(t)
		errorHandlerMock.EXPECT().HandleAppError(gin_context, middleware.ErrNoUID, http.StatusUnauthorized)
		getAll := GetAll(todoRepositoryMock, errorHandlerMock)
		getAll(gin_context)
	})
}

func TestGetById(t *testing.T) {
	t.Run("Good case", func(t *testing.T) {
		todoRepositoryMock, gin_context, http_recorder, errorHandlerMock := createMocks(t)
		todoId := uuid.New()
		token := &auth.Token{UID: "heowh"}
		uUidParseMock := func(id string) (uuid.UUID, error) {
			return todoId, nil
		}
		done := false
		ti, _ := time.Parse(time.RFC3339, "2022-07-21T14:07:05.768Z")
		todo := model.Todo{Id: todoId.String(), Title: "title1",
			Description: "description1", Done: &done, CreatedAt: ti}
		gin_context.Params = append(gin_context.Params, gin.Param{Key: "id", Value: todoId.String()})
		gin_context.Set(middleware.AuthToken, token)
		todoRepositoryMock.EXPECT().GetById(todoId.String(), token.UID).Return(&todo, nil)
		getById := GetById(todoRepositoryMock, errorHandlerMock, uUidParseMock)
		getById(gin_context)
		assert.Equal(t, http.StatusOK, http_recorder.Code)
		var got model.Todo
		err := json.Unmarshal(http_recorder.Body.Bytes(), &got)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, todo, got)
	})

	t.Run("When invalid id is sent as a path parameter in the url", func(t *testing.T) {
		todoRepositoryMock, gin_context, _, errorHandlerMock := createMocks(t)
		token := &auth.Token{UID: "heowh"}
		uUidParseMock := func(id string) (uuid.UUID, error) {
			return uuid.Nil, common.ErrError
		}
		gin_context.Params = append(gin_context.Params, gin.Param{Key: "id", Value: "oehwegiuf"})
		gin_context.Set(middleware.AuthToken, token)
		todoRepositoryMock.EXPECT().GetById(gomock.Any(), token.UID).Times(0)
		errorHandlerMock.EXPECT().HandleAppError(gin_context, common.ErrError, http.StatusBadRequest)
		getById := GetById(todoRepositoryMock, errorHandlerMock, uUidParseMock)
		getById(gin_context)
	})

	t.Run("When TodoRepository returns an Error", func(t *testing.T) {
		todoRepositoryMock, gin_context, _, errorHandlerMock := createMocks(t)
		todoId := uuid.New()
		token := &auth.Token{UID: "heowh"}
		uUidParseMock := func(id string) (uuid.UUID, error) {
			return todoId, nil
		}
		gin_context.Params = append(gin_context.Params, gin.Param{Key: "id", Value: "oehwegiuf"})
		gin_context.Set(middleware.AuthToken, token)
		todoRepositoryMock.EXPECT().GetById(gomock.Any(), token.UID).Return(nil, common.ErrError)
		errorHandlerMock.EXPECT().HandleAppError(gin_context, common.ErrError, http.StatusInternalServerError)
		getById := GetById(todoRepositoryMock, errorHandlerMock, uUidParseMock)
		getById(gin_context)
	})

	t.Run("When parse is nil", func(t *testing.T) {
		todoRepositoryMock, gin_context, _, errorHandlerMock := createMocks(t)
		token := &auth.Token{UID: "heowh"}
		gin_context.Set(middleware.AuthToken, token)
		errorHandlerMock.EXPECT().HandleAppError(gin_context, ErrParseIsNil, http.StatusInternalServerError)
		getById := GetById(todoRepositoryMock, errorHandlerMock, nil)
		getById(gin_context)
	})

	t.Run("When there is no auth token in the web context", func(t *testing.T) {
		todoRepositoryMock, gin_context, _, errorHandlerMock := createMocks(t)
		errorHandlerMock.EXPECT().HandleAppError(gin_context, middleware.ErrNoUID, http.StatusUnauthorized)
		getById := GetById(todoRepositoryMock, errorHandlerMock, nil)
		getById(gin_context)
	})
}

func TestUpdate(t *testing.T) {
	t.Run("Good case", func(t *testing.T) {
		todoRepositoryMock, gin_context, http_recorder, errorHandlerMock := createMocks(t)
		update := Update(todoRepositoryMock, errorHandlerMock)
		done := false
		token := &auth.Token{UID: "nfwseo"}
		ti, _ := time.Parse(time.RFC3339, "2022-09-21T14:07:05.768Z")
		todo := model.Todo{Id: uuid.New().String(), Title: "title1",
			Description: "description1",
			Done:        &done, CreatedAt: ti}
		json_bytes, err := json.Marshal(todo)
		if err != nil {
			t.Fatal(err)
		}
		web_request := &http.Request{
			Body:   io.NopCloser(bytes.NewBuffer(json_bytes)),
			Header: map[string][]string{"Content-Type": {"application/json"}}}
		gin_context.Request = web_request
		gin_context.Set(middleware.AuthToken, token)
		todoRepositoryMock.EXPECT().Update(&todo, token.UID).Return(nil)
		update(gin_context)
		assert.Equal(t, http.StatusNoContent, http_recorder.Code)
		assert.Empty(t, http_recorder.Body.Bytes())
	})

	t.Run("When required fields are not present in the web request body", func(t *testing.T) {
		todoRepositoryMock, gin_context, _, errorHandlerMock := createMocks(t)
		done := false
		token := &auth.Token{UID: "nfwseo"}
		ti, _ := time.Parse(time.RFC3339, "2022-09-21T14:07:05.768Z")
		todo := model.Todo{Id: uuid.New().String(), Title: "title1",
			Done: &done, CreatedAt: ti}
		json_bytes, err := json.Marshal(todo)
		if err != nil {
			t.Fatal(err)
		}
		web_request := &http.Request{
			Body:   io.NopCloser(bytes.NewBuffer(json_bytes)),
			Header: map[string][]string{"Content-Type": {"application/json"}}}
		gin_context.Request = web_request
		gin_context.Set(middleware.AuthToken, token)
		update := Update(todoRepositoryMock, errorHandlerMock)
		todoRepositoryMock.EXPECT().Update(gomock.Any(), token.UID).Times(0)
		errorHandlerMock.EXPECT().HandleAppError(gin_context, gomock.Any(), http.StatusBadRequest).
			DoAndReturn(func(ctx *gin.Context, err error, code int) {
				if !strings.Contains(err.Error(), "Description") {
					t.Fail()
				}
			})
		update(gin_context)
	})

	t.Run("When TodoRepository returns an error", func(t *testing.T) {
		todoRepositoryMock, gin_context, _, errorHandlerMock := createMocks(t)
		update := Update(todoRepositoryMock, errorHandlerMock)
		done := false
		token := &auth.Token{UID: "nfwseo"}
		ti, _ := time.Parse(time.RFC3339, "2022-09-21T14:07:05.768Z")
		todo := model.Todo{Id: uuid.New().String(), Title: "title1",
			Description: "description1",
			Done:        &done, CreatedAt: ti}
		json_bytes, err := json.Marshal(todo)
		if err != nil {
			t.Fatal(err)
		}
		web_request := &http.Request{
			Body:   io.NopCloser(bytes.NewBuffer(json_bytes)),
			Header: map[string][]string{"Content-Type": {"application/json"}}}
		gin_context.Request = web_request
		gin_context.Set(middleware.AuthToken, token)
		todoRepositoryMock.EXPECT().Update(&todo, token.UID).Return(common.ErrError)
		errorHandlerMock.EXPECT().HandleAppError(gin_context, common.ErrError, http.StatusInternalServerError)
		update(gin_context)
	})

	t.Run("When there is no auth token in the web context", func(t *testing.T) {
		todoRepositoryMock, gin_context, _, errorHandlerMock := createMocks(t)
		update := Update(todoRepositoryMock, errorHandlerMock)
		done := false
		ti, _ := time.Parse(time.RFC3339, "2022-09-21T14:07:05.768Z")
		todo := model.Todo{Id: uuid.New().String(), Title: "title1",
			Description: "description1",
			Done:        &done, CreatedAt: ti}
		json_bytes, err := json.Marshal(todo)
		if err != nil {
			t.Fatal(err)
		}
		web_request := &http.Request{
			Body:   io.NopCloser(bytes.NewBuffer(json_bytes)),
			Header: map[string][]string{"Content-Type": {"application/json"}}}
		gin_context.Request = web_request
		errorHandlerMock.EXPECT().HandleAppError(gin_context, middleware.ErrNoUID, http.StatusUnauthorized)
		update(gin_context)
	})
}

func TestDelete(t *testing.T) {
	t.Run("Good case", func(t *testing.T) {
		todoRepositoryMock, gin_context, http_recorder, errorHandlerMock := createMocks(t)
		token := &auth.Token{UID: "oiwhbegfwh"}
		todoId := uuid.New()
		uUidParseMock := func(id string) (uuid.UUID, error) {
			return todoId, nil
		}
		gin_context.Params = append(gin_context.Params, gin.Param{Key: "id", Value: todoId.String()})
		gin_context.Set(middleware.AuthToken, token)
		todoRepositoryMock.EXPECT().Delete(todoId.String(), token.UID).Return(nil)
		delete := Delete(todoRepositoryMock, errorHandlerMock, uUidParseMock)
		delete(gin_context)
		assert.Equal(t, http.StatusNoContent, http_recorder.Code)
		assert.Empty(t, http_recorder.Body.Bytes())
	})

	t.Run("When invalid todo id is sent as a path parameter in the url", func(t *testing.T) {
		todoRepositoryMock, gin_context, _, errorHandlerMock := createMocks(t)
		token := &auth.Token{UID: "oiwhbegfwh"}
		uUidParseMock := func(id string) (uuid.UUID, error) {
			return uuid.Nil, common.ErrError
		}
		gin_context.Params = append(gin_context.Params, gin.Param{Key: "id", Value: "oehwegiuf"})
		gin_context.Set(middleware.AuthToken, token)
		todoRepositoryMock.EXPECT().Delete(gomock.Any(), token.UID).Times(0)
		errorHandlerMock.EXPECT().HandleAppError(gin_context, common.ErrError, http.StatusBadRequest)
		delete := Delete(todoRepositoryMock, errorHandlerMock, uUidParseMock)
		delete(gin_context)
	})

	t.Run("When TodoRepository return an error", func(t *testing.T) {
		todoRepositoryMock, gin_context, _, errorHandlerMock := createMocks(t)
		token := &auth.Token{UID: "oiwhbegfwh"}
		todoId := uuid.New()
		uUidParseMock := func(id string) (uuid.UUID, error) {
			return todoId, nil
		}
		gin_context.Params = append(gin_context.Params, gin.Param{Key: "id", Value: todoId.String()})
		gin_context.Set(middleware.AuthToken, token)
		todoRepositoryMock.EXPECT().Delete(todoId.String(), token.UID).Return(common.ErrError)
		errorHandlerMock.EXPECT().HandleAppError(gin_context, common.ErrError, http.StatusInternalServerError)
		delete := Delete(todoRepositoryMock, errorHandlerMock, uUidParseMock)
		delete(gin_context)
	})

	t.Run("When parse is nil", func(t *testing.T) {
		todoRepositoryMock, gin_context, _, errorHandlerMock := createMocks(t)
		token := &auth.Token{UID: "oiwhbegfwh"}
		gin_context.Set(middleware.AuthToken, token)
		errorHandlerMock.EXPECT().HandleAppError(gin_context, ErrParseIsNil, http.StatusInternalServerError)
		delete := Delete(todoRepositoryMock, errorHandlerMock, nil)
		delete(gin_context)
	})

	t.Run("When there is no auth token in the web context", func(t *testing.T) {
		todoRepositoryMock, gin_context, _, errorHandlerMock := createMocks(t)
		errorHandlerMock.EXPECT().HandleAppError(gin_context, middleware.ErrNoUID, http.StatusUnauthorized)
		delete := Delete(todoRepositoryMock, errorHandlerMock, nil)
		delete(gin_context)
	})
}

func createMocks(t *testing.T) (*common.MockTodoRepository, *gin.Context, *httptest.ResponseRecorder, *common.MockErrorHandler) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	mockCtrl := gomock.NewController(t)
	http_recorder := httptest.NewRecorder()
	gin_context, _ := gin.CreateTestContext(http_recorder)
	return common.NewMockTodoRepository(mockCtrl), gin_context, http_recorder, common.NewMockErrorHandler(mockCtrl)
}
