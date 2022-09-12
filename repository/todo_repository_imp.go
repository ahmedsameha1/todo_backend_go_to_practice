package repository

import (
	"context"
	"errors"

	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/model"
	"github.com/google/uuid"
)

var ErrTodoRepositoryInitialization error = errors.New("DBPool is nil")
var ErrNotFound = errors.New("item is not found")
var ErrInvalidTodo = errors.New("invalid todo")

const (
	insertTodoQuery    string = "insert into todo (id, title, description, done, created_at, userId) values ($1::UUID, $2, $3, $4, $5, $6::UUID)"
	allTodosQuery      string = "select * from todo where user_id = $1::UUID"
	specificTodoQuery  string = "select * from todo where id = $1"
	updateQuery        string = "update todo set title = $2, description = $3, done = $4 where id = $1"
	deleteQuery        string = "delete from todo where id = $1"
)

type TodoRepositoryImpl struct {
	DBPool common.DBPool
}

func (tr TodoRepositoryImpl) Create(todo *model.Todo, userId string) (err error) {
	if !model.IsValid(todo) {
		return ErrInvalidTodo
	}
	if tr.DBPool == nil {
		err = ErrTodoRepositoryInitialization
	} else {
		_, err = tr.DBPool.Exec(context.Background(), insertTodoQuery, todo.Id, todo.Title,
			todo.Description, todo.Done, todo.CreatedAt, userId)
	}
	return err
}

func (tr TodoRepositoryImpl) GetAll(userId string) ([]model.Todo, error) {
	if tr.DBPool == nil {
		return nil, ErrTodoRepositoryInitialization
	}
	rows, _ := tr.DBPool.Query(context.Background(), allTodosQuery, userId)

	todos := []model.Todo{}
	for rows.Next() {
		var todo model.Todo
		if err := rows.Scan(&todo.Id); err != nil {
			return nil, err
		}
		if err := rows.Scan(&todo.Title); err != nil {
			return nil, err
		}
		if err := rows.Scan(&todo.Description); err != nil {
			return nil, err
		}
		if err := rows.Scan(&todo.Done); err != nil {
			return nil, err
		}
		if err := rows.Scan(&todo.CreatedAt); err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return todos, nil
}

func (tr TodoRepositoryImpl) GetById(id uuid.UUID) (*model.Todo, error) {
	if tr.DBPool == nil {
		return nil, ErrTodoRepositoryInitialization
	}
	rows, _ := tr.DBPool.Query(context.Background(), specificTodoQuery, id)
	if err := rows.Err(); err != nil {
		return nil, err
	} else {
		var todo model.Todo
		if rows.Next() {
			if err := rows.Scan(&todo.Id); err != nil {
				return nil, err
			}
			if err := rows.Scan(&todo.Title); err != nil {
				return nil, err
			}
			if err := rows.Scan(&todo.Description); err != nil {
				return nil, err
			}
			if err := rows.Scan(&todo.Done); err != nil {
				return nil, err
			}
			if err := rows.Scan(&todo.CreatedAt); err != nil {
				return nil, err
			}
			return &todo, nil
		} else {
			return nil, ErrNotFound
		}
	}
}

func (tr TodoRepositoryImpl) Update(todo *model.Todo) error {
	if !model.IsValid(todo) {
		return ErrInvalidTodo
	}
	if tr.DBPool == nil {
		return ErrTodoRepositoryInitialization
	}
	_, err := tr.DBPool.Exec(context.Background(), updateQuery, todo.Id, todo.Title, todo.Description, todo.Done)
	return err
}

func (tr TodoRepositoryImpl) Delete(id uuid.UUID) error {
	if tr.DBPool == nil {
		return ErrTodoRepositoryInitialization
	}
	_, err := tr.DBPool.Exec(context.Background(), deleteQuery, id)
	return err
}
