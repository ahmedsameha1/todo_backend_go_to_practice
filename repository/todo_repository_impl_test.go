package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/model"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestGetTodoRepository(t *testing.T) {
	t.Run("DBPool is nil", func(t *testing.T) {
		todoRepositroy, err := GetTodoRepository(nil)
		assert.Equal(t, ErrDBPoolIsNil, err)
		assert.Nil(t, todoRepositroy)
	})
	t.Run("DBPool is not nil", func(t *testing.T) {
		dbPool, _, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		todoRepositroy, err := GetTodoRepository(dbPool)
		assert.NotNil(t, todoRepositroy)
		assert.Nil(t, err)
	})
}

func TestCreate(t *testing.T) {
	t.Run("Good case", func(t *testing.T) {
		todoRepository, mock := create(t)
		todoDone := false
		ti, _ := time.Parse(time.RFC3339, "2022-09-21T14:07:05.768Z")
		userId := uuid.New().String()
		todo := model.Todo{Id: uuid.New().String(), Title: "title1",
			Description: "description1", Done: &todoDone, CreatedAt: ti}
		mock.ExpectExec(insertTodoQuery).WithArgs(todo.Id, todo.Title,
			todo.Description, todo.Done, todo.CreatedAt, userId).WillReturnResult(sqlmock.NewErrorResult(nil))
		err := todoRepository.Create(&todo, userId)
		assert.NoError(t, err)
		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("When DBPool returns an error", func(t *testing.T) {
		todoRepository, mock := create(t)
		todoDone := false
		ti, _ := time.Parse(time.RFC3339, "2022-09-21T14:07:05.768Z")
		userId := uuid.New().String()
		todo := model.Todo{Id: uuid.New().String(), Title: "title1",
			Description: "description1", Done: &todoDone, CreatedAt: ti}
		mock.ExpectExec(insertTodoQuery).WithArgs(todo.Id,
			todo.Title, todo.Description, todo.Done, todo.CreatedAt, userId).
			WillReturnError(common.ErrError)
		err := todoRepository.Create(&todo, userId)
		assert.Equal(t, common.ErrError, err)
		err = mock.ExpectationsWereMet()
		if err != nil {
			t.Error(err)
		}
	})

	t.Run("Invalid todo", func(t *testing.T) {
		todoRepository, _ := create(t)
		todoDone := false
		ti, _ := time.Parse(time.RFC3339, "2022-09-21T14:07:05.768Z")
		userId := uuid.New().String()
		invalidTodo := model.Todo{Title: "title1", Done: &todoDone, CreatedAt: ti}
		err := todoRepository.Create(&invalidTodo, userId)
		assert.Equal(t, ErrInvalidTodo, err)
	})

	t.Run("Invalid todo 2", func(t *testing.T) {
		todoRepository, _ := create(t)
		userId := uuid.New().String()
		err := todoRepository.Create(nil, userId)
		assert.Equal(t, ErrInvalidTodo, err)
	})
}

/*
func TestGetAll(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	dbRowsMock := common.NewMockDBRows(mockCtrl)
	userId := uuid.New().String()
	todoDone1 := false
	todoDone2 := true
	todoDone3 := false
	wantedTodos := []model.Todo{
		{Id: uuid.New().String(), Title: "title1", Description: "description1", Done: &todoDone1, CreatedAt: time.Now()},
		{Id: uuid.New().String(), Title: "title2", Description: "description2", Done: &todoDone2, CreatedAt: time.Now()},
		{Id: uuid.New().String(), Title: "title3", Description: "description3", Done: &todoDone3, CreatedAt: time.Now()},
	}
	todoRepositoryImpl = TodoRepositoryImpl{DBPool: dbPoolMock}
	gomock.InOrder(
		dbPoolMock.EXPECT().Query(gomock.Any(), allTodosQuery, userId).Return(dbRowsMock, nil),
		dbRowsMock.EXPECT().Next().Return(true),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[0].Id),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[0].Title),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[0].Description),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[0].Done),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[0].CreatedAt),
		dbRowsMock.EXPECT().Next().Return(true),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[1].Id),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[1].Title),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[1].Description),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[1].Done),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[1].CreatedAt),
		dbRowsMock.EXPECT().Next().Return(true),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[2].Id),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[2].Title),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[2].Description),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[2].Done),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[2].CreatedAt),
		dbRowsMock.EXPECT().Next().Return(false),
		dbRowsMock.EXPECT().Err().Return(nil),
	)
	todos, err := todoRepositoryImpl.GetAll(userId)
	assert.Equal(t, wantedTodos, todos)
	assert.NoError(t, err)
}

func TestGetAll2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	dbRowsMock := common.NewMockDBRows(mockCtrl)
	userId := uuid.New().String()
	todoRepositoryImpl = TodoRepositoryImpl{DBPool: dbPoolMock}
	gomock.InOrder(
		dbPoolMock.EXPECT().Query(gomock.Any(), allTodosQuery, userId).Return(dbRowsMock, nil),
		dbRowsMock.EXPECT().Next().Return(false),
		dbRowsMock.EXPECT().Err().Return(nil),
	)
	todos, err := todoRepositoryImpl.GetAll(userId)
	assert.Equal(t, []model.Todo{}, todos)
	assert.NoError(t, err)
}

func TestGetAllWhenAScanCallReturnAnError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	dbRowsMock := common.NewMockDBRows(mockCtrl)
	userId := uuid.New().String()
	todoDone1 := false
	todoDone2 := true
	todoDone3 := false
	wantedTodos := []model.Todo{
		{Id: uuid.New().String(), Title: "title1", Description: "description1", Done: &todoDone1, CreatedAt: time.Now()},
		{Id: uuid.New().String(), Title: "title2", Description: "description2", Done: &todoDone2, CreatedAt: time.Now()},
		{Id: uuid.New().String(), Title: "title3", Description: "description3", Done: &todoDone3, CreatedAt: time.Now()},
	}
	todoRepositoryImpl = TodoRepositoryImpl{DBPool: dbPoolMock}
	gomock.InOrder(
		dbPoolMock.EXPECT().Query(gomock.Any(), allTodosQuery, userId).Return(dbRowsMock, nil),
		dbRowsMock.EXPECT().Next().Return(true),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[0].Id),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[0].Title).Return(common.ErrError),
	)
	todos, err := todoRepositoryImpl.GetAll(userId)
	assert.Nil(t, todos)
	assert.Error(t, err, common.ErrError)
}

func TestGetAllWhenAScanCallReturnAnError2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	dbRowsMock := common.NewMockDBRows(mockCtrl)
	userId := uuid.New().String()
	todoDone1 := false
	todoDone2 := true
	todoDone3 := false
	wantedTodos := []model.Todo{
		{Id: uuid.New().String(), Title: "title1", Description: "description1", Done: &todoDone1, CreatedAt: time.Now()},
		{Id: uuid.New().String(), Title: "title2", Description: "description2", Done: &todoDone2, CreatedAt: time.Now()},
		{Id: uuid.New().String(), Title: "title3", Description: "description3", Done: &todoDone3, CreatedAt: time.Now()},
	}
	todoRepositoryImpl = TodoRepositoryImpl{DBPool: dbPoolMock}
	gomock.InOrder(
		dbPoolMock.EXPECT().Query(gomock.Any(), allTodosQuery, userId).Return(dbRowsMock, nil),
		dbRowsMock.EXPECT().Next().Return(true),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[0].Id),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[0].Title),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[0].Description),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[0].Done),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[0].CreatedAt),
		dbRowsMock.EXPECT().Next().Return(true),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[1].Id),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[1].Title),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[1].Description),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[1].Done),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[1].CreatedAt),
		dbRowsMock.EXPECT().Next().Return(true),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[2].Id),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[2].Title),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[2].Description),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[2].Done),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[2].CreatedAt).Return(common.ErrError),
	)
	todos, err := todoRepositoryImpl.GetAll(userId)
	assert.Nil(t, todos)
	assert.Error(t, common.ErrError, err)
}

func TestGetAllWhenErrCallReturnAnError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	dbRowsMock := common.NewMockDBRows(mockCtrl)
	userId := uuid.New().String()
	todoDone1 := false
	todoDone2 := true
	todoDone3 := false
	wantedTodos := []model.Todo{
		{Id: uuid.New().String(), Title: "title1", Description: "description1", Done: &todoDone1, CreatedAt: time.Now()},
		{Id: uuid.New().String(), Title: "title2", Description: "description2", Done: &todoDone2, CreatedAt: time.Now()},
		{Id: uuid.New().String(), Title: "title3", Description: "description3", Done: &todoDone3, CreatedAt: time.Now()},
	}
	todoRepositoryImpl = TodoRepositoryImpl{DBPool: dbPoolMock}
	gomock.InOrder(
		dbPoolMock.EXPECT().Query(gomock.Any(), allTodosQuery, userId).Return(dbRowsMock, common.ErrError),
		dbRowsMock.EXPECT().Next().Return(true),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[0].Id),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[0].Title),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[0].Description),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[0].Done),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[0].CreatedAt),
		dbRowsMock.EXPECT().Next().Return(true),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[1].Id),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[1].Title),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[1].Description),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[1].Done),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[1].CreatedAt),
		dbRowsMock.EXPECT().Next().Return(true),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[2].Id),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[2].Title),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[2].Description),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[2].Done),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[2].CreatedAt),
		dbRowsMock.EXPECT().Next().Return(false),
		dbRowsMock.EXPECT().Err().Return(common.ErrError),
	)
	todos, err := todoRepositoryImpl.GetAll(userId)
	assert.Nil(t, todos)
	assert.Error(t, err, common.ErrError)
}

func TestGetAllWhenDBPoolIsNil(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	dbRowsMock := common.NewMockDBRows(mockCtrl)
	todoRepositoryImpl = TodoRepositoryImpl{DBPool: nil}
	userId := uuid.New().String()
	dbPoolMock.EXPECT().Query(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	dbRowsMock.EXPECT().Next().Times(0)
	dbRowsMock.EXPECT().Scan(gomock.Any()).Times(0)
	dbRowsMock.EXPECT().Err().Times(0)
	todos, err := todoRepositoryImpl.GetAll(userId)
	assert.Nil(t, todos)
	assert.Equal(t, err, ErrTodoRepositoryInitialization)
}

func TestGetById(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	dbRowsMock := common.NewMockDBRows(mockCtrl)
	todoRepositoryImpl = TodoRepositoryImpl{DBPool: dbPoolMock}
	userId := uuid.New().String()
	todoId := uuid.New().String()
	todoDone := false
	wantedTodo := model.Todo{Id: todoId, Title: "title1",
		Description: "description1", Done: &todoDone, CreatedAt: time.Now()}
	gomock.InOrder(
		dbPoolMock.EXPECT().Query(gomock.Any(), specificTodoQuery, todoId, userId).Return(dbRowsMock, nil),
		dbRowsMock.EXPECT().Err().Return(nil),
		dbRowsMock.EXPECT().Next().Return(true),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodo.Id),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodo.Title),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodo.Description),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodo.Done),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodo.CreatedAt),
	)
	todo, err := todoRepositoryImpl.GetById(todoId, userId)
	assert.Equal(t, wantedTodo, *todo)
	assert.Nil(t, err)
}

func TestGetByIdWhenNotFound(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	dbRowsMock := common.NewMockDBRows(mockCtrl)
	todoRepositoryImpl = TodoRepositoryImpl{DBPool: dbPoolMock}
	userId := uuid.New().String()
	todoId := uuid.New().String()
	gomock.InOrder(
		dbPoolMock.EXPECT().Query(gomock.Any(), specificTodoQuery, todoId, userId).Return(dbRowsMock, nil),
		dbRowsMock.EXPECT().Err().Return(nil),
		dbRowsMock.EXPECT().Next().Return(false),
	)
	todo, err := todoRepositoryImpl.GetById(todoId, userId)
	assert.Nil(t, todo)
	assert.Equal(t, ErrNotFound, err)
}

func TestGetByIdWhenScanReturnAnError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	dbRowsMock := common.NewMockDBRows(mockCtrl)
	todoRepositoryImpl = TodoRepositoryImpl{DBPool: dbPoolMock}
	userId := uuid.New().String()
	todoId := uuid.New().String()
	todoDone := false
	wantedTodo := model.Todo{Id: todoId, Title: "title1",
		Description: "description1", Done: &todoDone, CreatedAt: time.Now()}
	gomock.InOrder(
		dbPoolMock.EXPECT().Query(gomock.Any(), specificTodoQuery, todoId, userId).Return(dbRowsMock, nil),
		dbRowsMock.EXPECT().Err().Return(nil),
		dbRowsMock.EXPECT().Next().Return(true),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodo.Id),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodo.Title).Return(common.ErrError),
	)
	todo, err := todoRepositoryImpl.GetById(todoId, userId)
	assert.Nil(t, todo)
	assert.NotNil(t, err)
}

func TestGetByIdWhenScanReturnAnError2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	dbRowsMock := common.NewMockDBRows(mockCtrl)
	todoRepositoryImpl = TodoRepositoryImpl{DBPool: dbPoolMock}
	userId := uuid.New().String()
	todoId := uuid.New().String()
	todoDone := false
	wantedTodo := model.Todo{Id: todoId, Title: "title1",
		Description: "description1", Done: &todoDone, CreatedAt: time.Now()}
	gomock.InOrder(
		dbPoolMock.EXPECT().Query(gomock.Any(), specificTodoQuery, todoId, userId).Return(dbRowsMock, nil),
		dbRowsMock.EXPECT().Err().Return(nil),
		dbRowsMock.EXPECT().Next().Return(true),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodo.Id),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodo.Title),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodo.Description),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodo.Done),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodo.CreatedAt).Return(common.ErrError),
	)
	todo, err := todoRepositoryImpl.GetById(todoId, userId)
	assert.Nil(t, todo)
	assert.Equal(t, common.ErrError, err)
}

func TestGetByIdWhenErrReturnAnError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	dbRowsMock := common.NewMockDBRows(mockCtrl)
	todoRepositoryImpl = TodoRepositoryImpl{DBPool: dbPoolMock}
	userId := uuid.New().String()
	todoId := uuid.New().String()
	gomock.InOrder(
		dbPoolMock.EXPECT().Query(gomock.Any(), specificTodoQuery, todoId, userId).Return(dbRowsMock, nil),
		dbRowsMock.EXPECT().Err().Return(common.ErrError),
	)
	todo, err := todoRepositoryImpl.GetById(todoId, userId)
	assert.Nil(t, todo)
	assert.Equal(t, common.ErrError, err)
}

func TestGetByIdWhenDBPoolIsNil(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	dbRowsMock := common.NewMockDBRows(mockCtrl)
	todoRepositoryImpl = TodoRepositoryImpl{DBPool: nil}
	userId := uuid.New().String()
	todoId := uuid.New().String()
	dbPoolMock.EXPECT().Query(gomock.Any(), gomock.Any()).Times(0)
	dbRowsMock.EXPECT().Err().Times(0)
	dbRowsMock.EXPECT().Next().Times(0)
	dbRowsMock.EXPECT().Scan(gomock.Any()).Times(0)
	todo, err := todoRepositoryImpl.GetById(todoId, userId)
	assert.Nil(t, todo)
	assert.Equal(t, ErrTodoRepositoryInitialization, err)
}

func TestUpdate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	todoRepositoryImpl = TodoRepositoryImpl{DBPool: dbPoolMock}
	userId := uuid.New().String()
	todoDone1 := false
	todo := model.Todo{Id: uuid.New().String(), Title: "title1", Description: "description1",
		Done: &todoDone1, CreatedAt: time.Now()}
	dbPoolMock.EXPECT().Exec(gomock.Any(), updateQuery, todo.Id, todo.Title,
		todo.Description, todo.Done, todo.CreatedAt, userId).Return(nil, nil)
	err := todoRepositoryImpl.Update(&todo, userId)
	assert.NoError(t, err)
}

func TestUpdateWhenDBPoolReturnAnError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	todoRepositoryImpl = TodoRepositoryImpl{DBPool: dbPoolMock}
	userId := uuid.New().String()
	todoDone1 := false
	todo := model.Todo{Id: uuid.New().String(), Title: "title1", Description: "description1",
		Done: &todoDone1, CreatedAt: time.Now()}
	dbPoolMock.EXPECT().Exec(gomock.Any(), updateQuery, todo.Id, todo.Title,
		todo.Description, todo.Done, todo.CreatedAt, userId).Return(nil, common.ErrError)
	err := todoRepositoryImpl.Update(&todo, userId)
	assert.Equal(t, common.ErrError, err)
}

func TestUpdateWhenDBoolIsNil(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	todoRepositoryImpl = TodoRepositoryImpl{DBPool: nil}
	userId := uuid.New().String()
	todoDone1 := false
	todo := model.Todo{Id: uuid.New().String(), Title: "title1", Description: "description1",
		Done: &todoDone1, CreatedAt: time.Now()}
	dbPoolMock.EXPECT().Exec(gomock.Any(), updateQuery, todo.Id, todo.Title,
		todo.Description, todo.Done, todo.CreatedAt, userId).Times(0)
	err := todoRepositoryImpl.Update(&todo, userId)
	assert.Equal(t, ErrTodoRepositoryInitialization, err)
}

func TestUpdateWhenTodoIsInvalid(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	todoRepositoryImpl = TodoRepositoryImpl{DBPool: dbPoolMock}
	userId := uuid.New().String()
	todoDone1 := false
	invalidTodo := model.Todo{Id: uuid.New().String(), Title: "", Description: "description1",
		Done: &todoDone1, CreatedAt: time.Now()}
	dbPoolMock.EXPECT().Exec(gomock.Any(), updateQuery, invalidTodo.Id, invalidTodo.Title,
		invalidTodo.Description, invalidTodo.Done, invalidTodo.CreatedAt, userId).Times(0)
	err := todoRepositoryImpl.Update(&invalidTodo, userId)
	assert.Equal(t, ErrInvalidTodo, err)
}

func TestUpdateWhenTodoIsInvalid2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	todoRepositoryImpl = TodoRepositoryImpl{DBPool: dbPoolMock}
	userId := uuid.New().String()
	todoDone1 := false
	invalidTodo := model.Todo{Id: uuid.New().String(), Title: "", Description: "description1",
		Done: &todoDone1, CreatedAt: time.Now()}
	dbPoolMock.EXPECT().Exec(gomock.Any(), updateQuery, invalidTodo.Id, invalidTodo.Title,
		invalidTodo.Description, invalidTodo.Done, invalidTodo.CreatedAt, userId).Times(0)
	err := todoRepositoryImpl.Update(nil, userId)
	assert.Equal(t, ErrInvalidTodo, err)
}

func TestDelete(t *testing.T) {
	func1 := func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		dbPoolMock := common.NewMockDBPool(mockCtrl)
		todoRepositoryImpl = TodoRepositoryImpl{DBPool: dbPoolMock}
		userId := uuid.New().String()
		todoId := uuid.New().String()
		dbPoolMock.EXPECT().Exec(gomock.Any(), deleteQuery, todoId, userId).Return([]byte{}, nil)
		err := todoRepositoryImpl.Delete(todoId, userId)
		assert.NoError(t, err)
	}
	t.Run("When DBPool.Exec returns no error", func1)
	t.Run("When DBPool.Exec returns an error", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		dbPoolMock := common.NewMockDBPool(mockCtrl)
		todoRepositoryImpl = TodoRepositoryImpl{DBPool: dbPoolMock}
		userId := uuid.New().String()
		todoId := uuid.New().String()
		dbPoolMock.EXPECT().Exec(gomock.Any(), deleteQuery, todoId, userId).Return([]byte{}, common.ErrError)
		err := todoRepositoryImpl.Delete(todoId, userId)
		assert.Equal(t, common.ErrError, err)
	})
	t.Run("When DBPool is nil", func(t *testing.T) {
		mockCtrl := gomock.NewController(t)
		dbPoolMock := common.NewMockDBPool(mockCtrl)
		todoRepositoryImpl = TodoRepositoryImpl{DBPool: nil}
		userId := uuid.New().String()
		todoId := uuid.New().String()
		dbPoolMock.EXPECT().Exec(gomock.Any(), deleteQuery, todoId, userId).Times(0)
		err := todoRepositoryImpl.Delete(todoId, userId)
		assert.Equal(t, ErrTodoRepositoryInitialization, err)
	})
}
*/

func create(t *testing.T) (common.TodoRepository, sqlmock.Sqlmock) {
	t.Helper()
	dbPool, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatal()
	}
	todoRepository, err := GetTodoRepository(dbPool)
	if err != nil {
		t.Fatal()
	}
	return todoRepository, mock
}
