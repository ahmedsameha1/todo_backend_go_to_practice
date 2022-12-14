package repository

import (
	"database/sql"
	"errors"

	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/model"
)

var ErrNotFound = errors.New("item is not found")
var ErrInvalidTodo = errors.New("invalid todo")
var ErrDBPoolIsNil = errors.New("DBPool is nil")

const (
	insertTodoQuery   string = "insert into todo (id, title, description, done, created_at, user_id) values ($1::UUID, $2, $3, $4, $5::timestamptz, $6)"
	allTodosQuery     string = "select id, title, description, done, created_at from todo where user_id = $1 order by created_at desc"
	specificTodoQuery string = "select id, title, description, done, created_at from todo where id = $1::UUID and user_id = $2"
	updateQuery       string = "update todo set title = $2, description = $3, done = $4, created_at = $5 where id = $1::UUID and user_id = $6"
	deleteQuery       string = "delete from todo where id = $1::UUID and user_id = $2"
)

type todoRepositoryImpl struct {
	DBPool *sql.DB
}

func GetTodoRepository(dbPool *sql.DB) (common.TodoRepository, error) {
	if dbPool == nil {
		return nil, ErrDBPoolIsNil
	}
	return todoRepositoryImpl{DBPool: dbPool}, nil
}

func (tr todoRepositoryImpl) Create(todo *model.Todo, userId string) (err error) {
	if !model.IsValid(todo) {
		return ErrInvalidTodo
	}
	_, err = tr.DBPool.Exec(insertTodoQuery, todo.Id, todo.Title,
		todo.Description, todo.Done, todo.CreatedAt, userId)
	return err
}

func (tr todoRepositoryImpl) GetAll(userId string) ([]model.Todo, error) {
	rows, err := tr.DBPool.Query(allTodosQuery, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	todos := []model.Todo{}
	for rows.Next() {
		var todo model.Todo
		if err := rows.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Done, &todo.CreatedAt); err != nil {
			return nil, err
		}
		todo.CreatedAt = todo.CreatedAt.UTC()
		todos = append(todos, todo)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return todos, nil
}

func (tr todoRepositoryImpl) GetById(id string, userId string) (*model.Todo, error) {
	row := tr.DBPool.QueryRow(specificTodoQuery, id, userId)
	var todo model.Todo
	if err := row.Scan(&todo.Id, &todo.Title, &todo.Description, &todo.Done, &todo.CreatedAt); err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		} else {
			return nil, err
		}
	} else {
		todo.CreatedAt = todo.CreatedAt.UTC()
		return &todo, nil
	}

}

func (tr todoRepositoryImpl) Update(todo *model.Todo, userId string) error {
	if !model.IsValid(todo) {
		return ErrInvalidTodo
	}
	_, err := tr.DBPool.Exec(updateQuery, todo.Id, todo.Title,
		todo.Description, todo.Done, todo.CreatedAt, userId)
	return err
}

func (tr todoRepositoryImpl) Delete(id string, userId string) error {
	_, err := tr.DBPool.Exec(deleteQuery, id, userId)
	return err
}
