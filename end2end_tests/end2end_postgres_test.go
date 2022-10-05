package end2endtests

import (
	"context"
	"log"
	"testing"

	"firebase.google.com/go/v4"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/handler"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/router"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/repository"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
	_ "github.com/jackc/pgx/v5/stdlib"
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
}