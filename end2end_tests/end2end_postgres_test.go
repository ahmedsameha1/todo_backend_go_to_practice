package end2endtests

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"testing"

	"firebase.google.com/go/v4"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/handler"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/middleware"
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
}
