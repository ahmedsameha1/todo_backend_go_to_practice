package end2endtests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	"firebase.google.com/go/v4"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/handler"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/middleware"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/model"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/repository"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/router"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/option"
)

func TestIt(t *testing.T) {
	sa := option.WithCredentialsFile("/home/ahmed/flutter-chat-max-37d0b-firebase-adminsdk-pr146-08fa0a9626.json")
	apiKeyBytes, err := os.ReadFile("/home/ahmed/flutter-chat-max-firebase-api-key")
	apiKey := strings.TrimRight(string(apiKeyBytes), "\n")
	if err != nil {
		log.Fatalln(err)
	}
	app, err := firebase.NewApp(context.Background(), nil, sa)
	if err != nil {
		log.Fatalln(err)
	}
	authClient, err := app.Auth(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	errorHandler := handler.ErrorHandlerImpl{Logger: log.Default()}
	container, todoRepository := repository.SetupPostgres(t)
	defer container.Terminate(context.Background())
	router := router.SetTodoRoutes(gin.Default(), todoRepository, errorHandler, authClient)
	go router.Run()
	toGetIdTokenRequestBody := `{"email":"test1@test.com","password":"password","returnSecureToken":true}`
	toGetIdTokenRequestUrl := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s", apiKey)
	toGetIdTokenWebRequest, err := http.NewRequest("POST", toGetIdTokenRequestUrl, bytes.NewBuffer([]byte(toGetIdTokenRequestBody)))
	if err != nil {
		log.Fatalln(err)
	}
	toGetIdTokenWebRequest.Header.Set("Content-Type", "application/json")
	toGetIdTokenWebResponse, err := http.DefaultClient.Do(toGetIdTokenWebRequest)
	if err != nil {
		log.Fatalln(err)
	}
	defer toGetIdTokenWebResponse.Body.Close()
	toGetIdTokenResponseBody, err := io.ReadAll(toGetIdTokenWebResponse.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var firebaseResponse gin.H
	err = json.Unmarshal(toGetIdTokenResponseBody, &firebaseResponse)
	if err != nil {
		log.Fatalln(err)
	}
	idToken := firebaseResponse["idToken"].(string)
	toGetIdTokenRequestBody2 := `{"email":"test2@test.com","password":"password2","returnSecureToken":true}`
	toGetIdTokenRequestUrl2 := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s", apiKey)
	toGetIdTokenWebRequest2, err := http.NewRequest("POST", toGetIdTokenRequestUrl2, bytes.NewBuffer([]byte(toGetIdTokenRequestBody2)))
	if err != nil {
		log.Fatalln(err)
	}
	toGetIdTokenWebRequest2.Header.Set("Content-Type", "application/json")
	toGetIdTokenWebResponse2, err := http.DefaultClient.Do(toGetIdTokenWebRequest2)
	if err != nil {
		log.Fatalln(err)
	}
	defer toGetIdTokenWebResponse2.Body.Close()
	toGetIdTokenResponseBody2, err := io.ReadAll(toGetIdTokenWebResponse2.Body)
	if err != nil {
		log.Fatalln(err)
	}
	var firebaseResponse2 gin.H
	err = json.Unmarshal(toGetIdTokenResponseBody2, &firebaseResponse2)
	if err != nil {
		log.Fatalln(err)
	}
	idToken2 := firebaseResponse2["idToken"].(string)
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
	toBeDeletedTodoId := uuid.New().String()
	ti2, _ := time.Parse(time.RFC3339, "2020-09-21T14:07:05.768Z")
	toBeDeletedTodo := model.Todo{
		Id:          toBeDeletedTodoId,
		Title:       "title2",
		Description: "description2",
		Done:        &todoDone,
		CreatedAt:   ti2,
	}

	t.Run("GET method - /todos: no authorization header", func(t *testing.T) {
		res, err := http.Get("http://localhost:8080/todos")
		if err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var bodyString gin.H
		err = json.Unmarshal(body, &bodyString)
		if err != nil {
			log.Fatalln(err)
		}
		assert.Equal(t, gin.H{"error": middleware.ErrNoAuthorizationHeader.Error()}, bodyString)
	})

	t.Run(`GET method - "/todos/:id" : no authorization header`, func(t *testing.T) {
		res, err := http.Get("http://localhost:8080/todos/" + uuid.New().String())
		if err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var bodyString gin.H
		err = json.Unmarshal(body, &bodyString)
		if err != nil {
			log.Fatalln(err)
		}
		assert.Equal(t, gin.H{"error": middleware.ErrNoAuthorizationHeader.Error()}, bodyString)
	})

	t.Run("POST method - /todos: no authorization header", func(t *testing.T) {
		res, err := http.Post("http://localhost:8080/todos", "", nil)
		if err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var bodyString gin.H
		err = json.Unmarshal(body, &bodyString)
		if err != nil {
			log.Fatalln(err)
		}
		assert.Equal(t, gin.H{"error": middleware.ErrNoAuthorizationHeader.Error()}, bodyString)
	})

	t.Run("PUT method - /todos: no authorization header", func(t *testing.T) {
		request, err := http.NewRequest("PUT", "http://localhost:8080/todos", nil)
		if err != nil {
			log.Fatalln(err)
		}
		res, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var bodyString gin.H
		err = json.Unmarshal(body, &bodyString)
		if err != nil {
			log.Fatalln(err)
		}
		assert.Equal(t, gin.H{"error": middleware.ErrNoAuthorizationHeader.Error()}, bodyString)
	})

	t.Run(`Delete method - "/todos/:id" : no authorization header`, func(t *testing.T) {
		request, err := http.NewRequest("DELETE", "http://localhost:8080/todos/"+uuid.New().String(), nil)
		if err != nil {
			log.Fatalln(err)
		}
		res, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()
		assert.Equal(t, http.StatusUnauthorized, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var bodyString gin.H
		err = json.Unmarshal(body, &bodyString)
		if err != nil {
			log.Fatalln(err)
		}
		assert.Equal(t, gin.H{"error": middleware.ErrNoAuthorizationHeader.Error()}, bodyString)
	})

	t.Run("POST method - /todos: Good case", func(t *testing.T) {
		expectedTodoJson, err := json.Marshal(expectedTodo)
		if err != nil {
			log.Fatalln(err)
		}
		request, err := http.NewRequest("POST", "http://localhost:8080/todos", bytes.NewBuffer(expectedTodoJson))
		request.Header.Set(middleware.AUTHORIZATION, middleware.BEARER+idToken)
		if err != nil {
			log.Fatalln(err)
		}
		res, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var todo model.Todo
		err = json.Unmarshal(body, &todo)
		if err != nil {
			log.Fatalln(err)
		}
		assert.Equal(t, expectedTodo, todo)
		request, err = http.NewRequest("GET", "http://localhost:8080/todos", nil)
		request.Header.Set(middleware.AUTHORIZATION, middleware.BEARER+idToken)
		if err != nil {
			log.Fatalln(err)
		}
		res, err = http.DefaultClient.Do(request)
		if err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		body, err = io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var todos []model.Todo
		err = json.Unmarshal(body, &todos)
		if err != nil {
			log.Fatalln(err)
		}
		returnedTodo := todos[0]
		assert.Equal(t, expectedTodo, returnedTodo)
	})

	t.Run("POST method - /todos: invalid todo", func(t *testing.T) {
		todoDone := true
		ti, _ := time.Parse(time.RFC3339, "2022-09-21T14:07:05.768Z")
		todoId := uuid.New().String()
		expectedTodo := model.Todo{
			Id:          todoId,
			Description: "description1",
			Done:        &todoDone,
			CreatedAt:   ti,
		}
		expectedTodoJson, err := json.Marshal(expectedTodo)
		if err != nil {
			log.Fatalln(err)
		}
		request, err := http.NewRequest("POST", "http://localhost:8080/todos", bytes.NewBuffer(expectedTodoJson))
		request.Header.Set(middleware.AUTHORIZATION, middleware.BEARER+idToken)
		if err != nil {
			log.Fatalln(err)
		}
		res, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var got map[string]string
		err = json.Unmarshal(body, &got)
		if err != nil {
			log.Fatalln(err)
		}
		assert.True(t, strings.Contains(got["error"], "Title"))
		request, err = http.NewRequest("GET", "http://localhost:8080/todos", nil)
		request.Header.Set(middleware.AUTHORIZATION, middleware.BEARER+idToken)
		if err != nil {
			log.Fatalln(err)
		}
		res, err = http.DefaultClient.Do(request)
		if err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		body, err = io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var todos []model.Todo
		err = json.Unmarshal(body, &todos)
		if err != nil {
			log.Fatalln(err)
		}
		assert.Len(t, todos, 1)
	})
	
	t.Run("PUT method - /todos: Good case", func(t *testing.T) {
		todoDone := true
		ti, _ := time.Parse(time.RFC3339, "2022-09-21T14:07:05.768Z")
		toBeUpdatedTodo := model.Todo{
			Id:          todoId,
			Title:       "title1updated",
			Description: "description1updated",
			Done:        &todoDone,
			CreatedAt:   ti,
		}
		toBeUpdatedTodoJson, err := json.Marshal(toBeUpdatedTodo)
		if err != nil {
			log.Fatalln(err)
		}
		request, err := http.NewRequest("PUT", "http://localhost:8080/todos", bytes.NewBuffer(toBeUpdatedTodoJson))
		request.Header.Set(middleware.AUTHORIZATION, middleware.BEARER+idToken)
		if err != nil {
			log.Fatalln(err)
		}
		res, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()
		assert.Equal(t, http.StatusNoContent, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		assert.Empty(t, body)
		request, err = http.NewRequest("GET", "http://localhost:8080/todos", nil)
		request.Header.Set(middleware.AUTHORIZATION, middleware.BEARER+idToken)
		if err != nil {
			log.Fatalln(err)
		}
		res, err = http.DefaultClient.Do(request)
		if err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		body, err = io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var todos []model.Todo
		err = json.Unmarshal(body, &todos)
		if err != nil {
			log.Fatalln(err)
		}
		returnedTodo := todos[0]
		assert.Equal(t, toBeUpdatedTodo, returnedTodo)
		expectedTodoJson, err := json.Marshal(expectedTodo)
		if err != nil {
			log.Fatalln(err)
		}
		request, err = http.NewRequest("PUT", "http://localhost:8080/todos", bytes.NewBuffer(expectedTodoJson))
		request.Header.Set(middleware.AUTHORIZATION, middleware.BEARER+idToken)
		if err != nil {
			log.Fatalln(err)
		}
		_, err = http.DefaultClient.Do(request)
		if err != nil {
			log.Fatalln(err)
		}
	})

	t.Run("PUT method - /todos: invalid todo", func(t *testing.T) {
		todoDone := true
		ti, _ := time.Parse(time.RFC3339, "2022-09-21T14:07:05.768Z")
		toBeUpdatedTodo := model.Todo{
			Id:          todoId,
			Description: "description1updated",
			Done:        &todoDone,
			CreatedAt:   ti,
		}
		toBeUpdatedTodoJson, err := json.Marshal(toBeUpdatedTodo)
		if err != nil {
			log.Fatalln(err)
		}
		request, err := http.NewRequest("PUT", "http://localhost:8080/todos", bytes.NewBuffer(toBeUpdatedTodoJson))
		request.Header.Set(middleware.AUTHORIZATION, middleware.BEARER+idToken)
		if err != nil {
			log.Fatalln(err)
		}
		res, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var got map[string]string
		err = json.Unmarshal(body, &got)
		if err != nil {
			log.Fatalln(err)
		}
		assert.True(t, strings.Contains(got["error"], "Title"))
		request, err = http.NewRequest("GET", "http://localhost:8080/todos", nil)
		request.Header.Set(middleware.AUTHORIZATION, middleware.BEARER+idToken)
		if err != nil {
			log.Fatalln(err)
		}
		res, err = http.DefaultClient.Do(request)
		if err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		body, err = io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var todos []model.Todo
		err = json.Unmarshal(body, &todos)
		if err != nil {
			log.Fatalln(err)
		}
		assert.Len(t, todos, 1)
		returnedTodo := todos[0]
		assert.Equal(t, expectedTodo, returnedTodo)
	})

	t.Run("PUT method - /todos: diferent todo id", func(t *testing.T) {
		todoDone := true
		ti, _ := time.Parse(time.RFC3339, "2022-09-21T14:07:05.768Z")
		toBeUpdatedTodo := model.Todo{
			Id:          uuid.New().String(),
			Title:       "title1updated",
			Description: "description1updated",
			Done:        &todoDone,
			CreatedAt:   ti,
		}
		toBeUpdatedTodoJson, err := json.Marshal(toBeUpdatedTodo)
		if err != nil {
			log.Fatalln(err)
		}
		request, err := http.NewRequest("PUT", "http://localhost:8080/todos", bytes.NewBuffer(toBeUpdatedTodoJson))
		request.Header.Set(middleware.AUTHORIZATION, middleware.BEARER+idToken)
		if err != nil {
			log.Fatalln(err)
		}
		res, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()
		assert.Equal(t, http.StatusNoContent, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		assert.Empty(t, body)
		request, err = http.NewRequest("GET", "http://localhost:8080/todos", nil)
		request.Header.Set(middleware.AUTHORIZATION, middleware.BEARER+idToken)
		if err != nil {
			log.Fatalln(err)
		}
		res, err = http.DefaultClient.Do(request)
		if err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		body, err = io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var todos []model.Todo
		err = json.Unmarshal(body, &todos)
		if err != nil {
			log.Fatalln(err)
		}
		returnedTodo := todos[0]
		assert.Equal(t, expectedTodo, returnedTodo)
	})

	t.Run("PUT method - /todos: diferent user id", func(t *testing.T) {
		todoDone := true
		ti, _ := time.Parse(time.RFC3339, "2022-09-21T14:07:05.768Z")
		toBeUpdatedTodo := model.Todo{
			Id:          todoId,
			Title:       "title1updated",
			Description: "description1updated",
			Done:        &todoDone,
			CreatedAt:   ti,
		}
		toBeUpdatedTodoJson, err := json.Marshal(toBeUpdatedTodo)
		if err != nil {
			log.Fatalln(err)
		}
		request, err := http.NewRequest("PUT", "http://localhost:8080/todos", bytes.NewBuffer(toBeUpdatedTodoJson))
		request.Header.Set(middleware.AUTHORIZATION, middleware.BEARER+idToken2)
		if err != nil {
			log.Fatalln(err)
		}
		res, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()
		assert.Equal(t, http.StatusNoContent, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		assert.Empty(t, body)
		request, err = http.NewRequest("GET", "http://localhost:8080/todos", nil)
		request.Header.Set(middleware.AUTHORIZATION, middleware.BEARER+idToken)
		if err != nil {
			log.Fatalln(err)
		}
		res, err = http.DefaultClient.Do(request)
		if err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		body, err = io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var todos []model.Todo
		err = json.Unmarshal(body, &todos)
		if err != nil {
			log.Fatalln(err)
		}
		returnedTodo := todos[0]
		assert.Equal(t, expectedTodo, returnedTodo)
	})

	t.Run("DELETE method - /todos/:id: Good case", func(t *testing.T) {
		toBeDeletedTodoJson, err := json.Marshal(toBeDeletedTodo)
		if err != nil {
			log.Fatalln(err)
		}
		request, err := http.NewRequest("POST", "http://localhost:8080/todos", bytes.NewBuffer(toBeDeletedTodoJson))
		request.Header.Set(middleware.AUTHORIZATION, middleware.BEARER+idToken)
		if err != nil {
			log.Fatalln(err)
		}
		http.DefaultClient.Do(request)
		request, err = http.NewRequest("GET", "http://localhost:8080/todos", nil)
		request.Header.Set(middleware.AUTHORIZATION, middleware.BEARER+idToken)
		if err != nil {
			log.Fatalln(err)
		}
		res, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var todos []model.Todo
		err = json.Unmarshal(body, &todos)
		if err != nil {
			log.Fatalln(err)
		}
		assert.Len(t, todos, 2)
		returnedTodo := todos[1]
		assert.Equal(t, toBeDeletedTodo, returnedTodo)
		request, err = http.NewRequest("DELETE", "http://localhost:8080/todos/" + toBeDeletedTodoId, nil)
		request.Header.Set(middleware.AUTHORIZATION, middleware.BEARER+idToken)
		if err != nil {
			log.Fatalln(err)
		}
		res, err = http.DefaultClient.Do(request)
		if err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()
		assert.Equal(t, http.StatusNoContent, res.StatusCode)
		body, err = io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		assert.Empty(t, body)
		request, err = http.NewRequest("GET", "http://localhost:8080/todos", nil)
		request.Header.Set(middleware.AUTHORIZATION, middleware.BEARER+idToken)
		if err != nil {
			log.Fatalln(err)
		}
		res, err = http.DefaultClient.Do(request)
		if err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		body, err = io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var todos1 []model.Todo
		err = json.Unmarshal(body, &todos1)
		if err != nil {
			log.Fatalln(err)
		}
		assert.Len(t, todos1, 1)
		returnedTodo = todos1[0]
		assert.Equal(t, expectedTodo, returnedTodo)
	})

	t.Run("DELETE method - /todos/:id: invalid todo id", func(t *testing.T) {
		toBeDeletedTodoJson, err := json.Marshal(toBeDeletedTodo)
		if err != nil {
			log.Fatalln(err)
		}
		request, err := http.NewRequest("POST", "http://localhost:8080/todos", bytes.NewBuffer(toBeDeletedTodoJson))
		request.Header.Set(middleware.AUTHORIZATION, middleware.BEARER+idToken)
		if err != nil {
			log.Fatalln(err)
		}
		http.DefaultClient.Do(request)
		request, err = http.NewRequest("GET", "http://localhost:8080/todos", nil)
		request.Header.Set(middleware.AUTHORIZATION, middleware.BEARER+idToken)
		if err != nil {
			log.Fatalln(err)
		}
		res, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var todos []model.Todo
		err = json.Unmarshal(body, &todos)
		if err != nil {
			log.Fatalln(err)
		}
		assert.Len(t, todos, 2)
		returnedTodo := todos[1]
		assert.Equal(t, toBeDeletedTodo, returnedTodo)
		request, err = http.NewRequest("DELETE", "http://localhost:8080/todos/" + "turw", nil)
		request.Header.Set(middleware.AUTHORIZATION, middleware.BEARER+idToken)
		if err != nil {
			log.Fatalln(err)
		}
		res, err = http.DefaultClient.Do(request)
		if err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()
		assert.Equal(t, http.StatusBadRequest, res.StatusCode)
		body, err = io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var got gin.H
		err = json.Unmarshal(body, &got)
		if err != nil {
			log.Fatalln(err)
		}
		assert.True(t, strings.Contains(got["error"].(string), "invalid"))	
		assert.True(t, strings.Contains(got["error"].(string), "UUID"))	
	})

	t.Run("DELETE method - /todos/:id: todo id is diferent", func(t *testing.T) {
		toBeDeletedTodoJson, err := json.Marshal(toBeDeletedTodo)
		if err != nil {
			log.Fatalln(err)
		}
		request, err := http.NewRequest("POST", "http://localhost:8080/todos", bytes.NewBuffer(toBeDeletedTodoJson))
		request.Header.Set(middleware.AUTHORIZATION, middleware.BEARER+idToken)
		if err != nil {
			log.Fatalln(err)
		}
		http.DefaultClient.Do(request)
		request, err = http.NewRequest("DELETE", "http://localhost:8080/todos/" + uuid.New().String(), nil)
		request.Header.Set(middleware.AUTHORIZATION, middleware.BEARER+idToken)
		if err != nil {
			log.Fatalln(err)
		}
		res, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()
		assert.Equal(t, http.StatusNoContent, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		assert.Empty(t, body)
		request, err = http.NewRequest("GET", "http://localhost:8080/todos", nil)
		request.Header.Set(middleware.AUTHORIZATION, middleware.BEARER+idToken)
		if err != nil {
			log.Fatalln(err)
		}
		res, err = http.DefaultClient.Do(request)
		if err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		body, err = io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var todos1 []model.Todo
		err = json.Unmarshal(body, &todos1)
		if err != nil {
			log.Fatalln(err)
		}
		assert.Len(t, todos1, 2)
	})

	t.Run("DELETE method - /todos/:id: user id is diferent", func(t *testing.T) {
		toBeDeletedTodoJson, err := json.Marshal(toBeDeletedTodo)
		if err != nil {
			log.Fatalln(err)
		}
		request, err := http.NewRequest("POST", "http://localhost:8080/todos", bytes.NewBuffer(toBeDeletedTodoJson))
		request.Header.Set(middleware.AUTHORIZATION, middleware.BEARER+idToken)
		if err != nil {
			log.Fatalln(err)
		}
		http.DefaultClient.Do(request)
		request, err = http.NewRequest("DELETE", "http://localhost:8080/todos/" + toBeDeletedTodoId, nil)
		request.Header.Set(middleware.AUTHORIZATION, middleware.BEARER+idToken2)
		if err != nil {
			log.Fatalln(err)
		}
		res, err := http.DefaultClient.Do(request)
		if err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()
		assert.Equal(t, http.StatusNoContent, res.StatusCode)
		body, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		assert.Empty(t, body)
		request, err = http.NewRequest("GET", "http://localhost:8080/todos", nil)
		request.Header.Set(middleware.AUTHORIZATION, middleware.BEARER+idToken)
		if err != nil {
			log.Fatalln(err)
		}
		res, err = http.DefaultClient.Do(request)
		if err != nil {
			log.Fatalln(err)
		}
		defer res.Body.Close()
		assert.Equal(t, http.StatusOK, res.StatusCode)
		body, err = io.ReadAll(res.Body)
		if err != nil {
			log.Fatalln(err)
		}
		var todos1 []model.Todo
		err = json.Unmarshal(body, &todos1)
		if err != nil {
			log.Fatalln(err)
		}
		assert.Len(t, todos1, 2)
	})

}
