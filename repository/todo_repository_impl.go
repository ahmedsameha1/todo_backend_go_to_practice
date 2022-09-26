package repository

import (
	"context"
	"errors"

	"github.com/ahmedsameha1/todo_backend_go_to_practice/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("item is not found")
var ErrInvalidTodo = errors.New("invalid todo")

const (
	insertTodoQuery   string = "insert into todo (id, title, description, done, created_at, user_id) values ($1::UUID, $2, $3, $4, $5::timestamptz, $6::UUID)"
	allTodosQuery     string = "select id, title, description, done, created_at from todo where user_id = $1::UUID order by created_at desc"
	specificTodoQuery string = "select id, title, description, done, created_at from todo where id = $1::UUID and user_id = $2::UUID"
	updateQuery       string = "update todo set title = $2, description = $3, done = $4, created_at = $5 where id = $1::UUID and user_id = $6::UUID"
	deleteQuery       string = "delete from todo where id = $1::UUID and user_id = $2::UUID"
)

type TodoRepositoryImpl struct {
	DBPool pgxpool.Pool
}

func (tr TodoRepositoryImpl) Create(todo *model.Todo, userId string) (err error) {
	if !model.IsValid(todo) {
		return ErrInvalidTodo
	}
		_, err = tr.DBPool.Exec(context.Background(), insertTodoQuery, todo.Id, todo.Title,
			todo.Description, todo.Done, todo.CreatedAt, userId)
	return err
}

func (tr TodoRepositoryImpl) GetAll(userId string) ([]model.Todo, error) {
	rows, _ := tr.DBPool.Query(context.Background(), allTodosQuery, userId)
	todos := []model.Todo{}
	for rows.Next() {
		var todo model.Todo
		if err := rows.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Done, &todo.CreatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return todos, nil
}

func (tr TodoRepositoryImpl) GetById(id string, userId string) (*model.Todo, error) {
	rows, _ := tr.DBPool.Query(context.Background(), specificTodoQuery, id, userId)
	if err := rows.Err(); err != nil {
		return nil, err
	} else {
		var todo model.Todo
		if rows.Next() {
			if err := rows.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Done, &todo.CreatedAt); err != nil {
				return nil, err
			}
			return &todo, nil
		} else {
			return nil, ErrNotFound
		}
	}
}

func (tr TodoRepositoryImpl) Update(todo *model.Todo, userId string) error {
	if !model.IsValid(todo) {
		return ErrInvalidTodo
	}
	_, err := tr.DBPool.Exec(context.Background(), updateQuery, todo.Id, todo.Title,
		todo.Description, todo.Done, todo.CreatedAt, userId)
	return err
}

func (tr TodoRepositoryImpl) Delete(id string, userId string) error {
	_, err := tr.DBPool.Exec(context.Background(), deleteQuery, id, userId)
	return err
}
