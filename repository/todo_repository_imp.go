package repository

import (
	"context"
	"errors"
	"time"

	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/model"
	"github.com/google/uuid"
)

var ErrTodoRepositoryInitialization error = errors.New("DBPool is nil")
var ErrTodoIsNil error = errors.New("this todo is nil")
var ErrNotFound = errors.New("item is not found")

const (
	insertTodoQuery   string = "insert into todo (id, title, description, done, createdAt) values ($1, $2, $3, $4, $5)"
	allTodosQuery     string = "select * from todo"
	specificTodoQuery string = "select * from todo where id = $1"
)

type TodoRepositoryImpl struct {
	DBPool             common.DBPool
	IDGenerator        func() uuid.UUID
	CreatedAtGenerator func() time.Time
}

func (tr TodoRepositoryImpl) Create(todo *model.Todo) (err error) {
	if tr.DBPool == nil || tr.IDGenerator == nil || tr.CreatedAtGenerator == nil {
		err = ErrTodoRepositoryInitialization
	} else if todo == nil {
		err = ErrTodoIsNil
	} else {
		_, err = tr.DBPool.Exec(context.Background(), insertTodoQuery, tr.IDGenerator(), todo.Title,
			todo.Description, todo.Done, tr.CreatedAtGenerator())
	}
	return err
}

func (tr TodoRepositoryImpl) GetAll() ([]model.Todo, error) {
	if tr.DBPool == nil {
		return nil, ErrTodoRepositoryInitialization
	}
	rows, _ := tr.DBPool.Query(context.Background(), allTodosQuery)

	todos := []model.Todo{}
	for rows.Next() {
		todo := model.Todo{}
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
		todo := model.Todo{}
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

/*
func (tr TodoRepositoryImpl) GetAllByUserId(id uuid.UUID) ([]model.Todo, error)
func (tr TodoRepositoryImpl) Update(todo *model.Todo) error
func (tr TodoRepositoryImpl) Delete(id uuid.UUID) error
*/
