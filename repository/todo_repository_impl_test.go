package repository

import (
	"errors"
	"testing"
	"time"

	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/model"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var anErrString string = "An error"
var anError error = errors.New(anErrString)

func TestCreate(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	id := uuid.New()
	now := time.Now()
	idGeneratorMock := func() uuid.UUID {
		return id
	}
	createdAtGeneratorMock := func() time.Time {
		return now
	}
	todoRepositoryImpl := TodoRepositoryImpl{DBPool: dbPoolMock, IDGenerator: idGeneratorMock,
		CreatedAtGenerator: createdAtGeneratorMock}
	todoDone := false
	todo := model.Todo{Title: "title1", Description: "description1", Done: &todoDone}
	dbPoolMock.EXPECT().Exec(gomock.Any(), insertTodoQuery, idGeneratorMock(),
		todo.Title, todo.Description, todo.Done, createdAtGeneratorMock()).Return(nil, nil)
	err := todoRepositoryImpl.Create(&todo)
	assert.NoError(t, err)
}

func TestCreateWhenDBPoolReturnAnError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	id := uuid.New()
	now := time.Now()
	idGeneratorMock := func() uuid.UUID {
		return id
	}
	createdAtGeneratorMock := func() time.Time {
		return now
	}
	todoRepositoryImpl := TodoRepositoryImpl{DBPool: dbPoolMock, IDGenerator: idGeneratorMock,
		CreatedAtGenerator: createdAtGeneratorMock}
	todoDone := false
	todo := model.Todo{Title: "title1", Description: "description1", Done: &todoDone}
	dbPoolMock.EXPECT().Exec(gomock.Any(), insertTodoQuery, idGeneratorMock(),
		todo.Title, todo.Description, todo.Done, createdAtGeneratorMock()).Return(nil, anError)
	err := todoRepositoryImpl.Create(&todo)
	assert.Error(t, err)
}

func TestCreateWhenCreatingTodoRepositoryWithNilDBPool(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	id := uuid.New()
	now := time.Now()
	idGeneratorMock := func() uuid.UUID {
		return id
	}
	createdAtGeneratorMock := func() time.Time {
		return now
	}
	todoRepositoryImpl := TodoRepositoryImpl{DBPool: nil, IDGenerator: idGeneratorMock,
		CreatedAtGenerator: createdAtGeneratorMock}
	todoDone := false
	todo := model.Todo{Title: "title1", Description: "description1", Done: &todoDone}
	dbPoolMock.EXPECT().Exec(gomock.Any(), gomock.Any()).Times(0)
	err := todoRepositoryImpl.Create(&todo)
	assert.Equal(t, err, ErrTodoRepositoryInitialization)
}

func TestCreateWhenCreatingTodoRepositoryWithNilIDGenerator(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	now := time.Now()
	createdAtGeneratorMock := func() time.Time {
		return now
	}
	todoRepositoryImpl := TodoRepositoryImpl{DBPool: dbPoolMock, IDGenerator: nil,
		CreatedAtGenerator: createdAtGeneratorMock}
	todoDone := false
	todo := model.Todo{Title: "title1", Description: "description1", Done: &todoDone}
	dbPoolMock.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	err := todoRepositoryImpl.Create(&todo)
	assert.Equal(t, err, ErrTodoRepositoryInitialization)
}

func TestCreateWhenCreatingTodoRepositoryWithNilCreatedAtGenerator(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	id := uuid.New()
	idGeneratorMock := func() uuid.UUID {
		return id
	}
	todoRepositoryImpl := TodoRepositoryImpl{DBPool: dbPoolMock, IDGenerator: idGeneratorMock,
		CreatedAtGenerator: nil}
	todoDone := false
	todo := model.Todo{Title: "title1", Description: "description1", Done: &todoDone}
	dbPoolMock.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	err := todoRepositoryImpl.Create(&todo)
	assert.Equal(t, err, ErrTodoRepositoryInitialization)
}

func TestCreateWhenCalledWithNilTodo(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	id := uuid.New()
	now := time.Now()
	idGeneratorMock := func() uuid.UUID {
		return id
	}
	createdAtGeneratorMock := func() time.Time {
		return now
	}
	todoRepositoryImpl := TodoRepositoryImpl{DBPool: dbPoolMock, IDGenerator: idGeneratorMock,
		CreatedAtGenerator: createdAtGeneratorMock}
	dbPoolMock.EXPECT().Exec(gomock.Any(), gomock.Any(), gomock.Any()).Times(0)
	err := todoRepositoryImpl.Create(nil)
	assert.Equal(t, err, ErrTodoIsNil)
}

func TestGetAll(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	dbRowsMock := common.NewMockDBRows(mockCtrl)
	id := uuid.New()
	now := time.Now()
	idGeneratorMock := func() uuid.UUID {
		return id
	}
	createdAtGeneratorMock := func() time.Time {
		return now
	}
	todoDone1 := false
	todoDone2 := true
	todoDone3 := false
	wantedTodos := []model.Todo{
		{Id: uuid.New(), Title: "title1", Description: "description1", Done: &todoDone1, CreatedAt: time.Now()},
		{Id: uuid.New(), Title: "title2", Description: "description2", Done: &todoDone2, CreatedAt: time.Now()},
		{Id: uuid.New(), Title: "title3", Description: "description3", Done: &todoDone3, CreatedAt: time.Now()},
	}
	todoRepositoryImpl := TodoRepositoryImpl{DBPool: dbPoolMock,
		IDGenerator: idGeneratorMock, CreatedAtGenerator: createdAtGeneratorMock}
	gomock.InOrder(
		dbPoolMock.EXPECT().Query(gomock.Any(), queryAllTodos).Return(dbRowsMock, nil),
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
	todos, err := todoRepositoryImpl.GetAll()
	assert.Equal(t, wantedTodos, todos)
	assert.NoError(t, err)
}

func TestGetAll2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	dbRowsMock := common.NewMockDBRows(mockCtrl)
	id := uuid.New()
	now := time.Now()
	idGeneratorMock := func() uuid.UUID {
		return id
	}
	createdAtGeneratorMock := func() time.Time {
		return now
	}
	todoRepositoryImpl := TodoRepositoryImpl{DBPool: dbPoolMock,
		IDGenerator: idGeneratorMock, CreatedAtGenerator: createdAtGeneratorMock}
	gomock.InOrder(
		dbPoolMock.EXPECT().Query(gomock.Any(), queryAllTodos).Return(dbRowsMock, nil),
		dbRowsMock.EXPECT().Next().Return(false),
		dbRowsMock.EXPECT().Err().Return(nil),
	)
	todos, err := todoRepositoryImpl.GetAll()
	assert.Equal(t, []model.Todo{}, todos)
	assert.NoError(t, err)
}

func TestGetAllWhenAScanCallReturnAnError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	dbRowsMock := common.NewMockDBRows(mockCtrl)
	id := uuid.New()
	now := time.Now()
	idGeneratorMock := func() uuid.UUID {
		return id
	}
	createdAtGeneratorMock := func() time.Time {
		return now
	}
	todoDone1 := false
	todoDone2 := true
	todoDone3 := false
	wantedTodos := []model.Todo{
		{Id: uuid.New(), Title: "title1", Description: "description1", Done: &todoDone1, CreatedAt: time.Now()},
		{Id: uuid.New(), Title: "title2", Description: "description2", Done: &todoDone2, CreatedAt: time.Now()},
		{Id: uuid.New(), Title: "title3", Description: "description3", Done: &todoDone3, CreatedAt: time.Now()},
	}
	todoRepositoryImpl := TodoRepositoryImpl{DBPool: dbPoolMock,
		IDGenerator: idGeneratorMock, CreatedAtGenerator: createdAtGeneratorMock}
	gomock.InOrder(
		dbPoolMock.EXPECT().Query(gomock.Any(), queryAllTodos).Return(dbRowsMock, nil),
		dbRowsMock.EXPECT().Next().Return(true),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[0].Id),
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[0].Title).Return(anError),
	)
	todos, err := todoRepositoryImpl.GetAll()
	assert.Nil(t, todos)
	assert.Error(t, err, anError)
}

func TestGetAllWhenAScanCallReturnAnError2(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	dbRowsMock := common.NewMockDBRows(mockCtrl)
	id := uuid.New()
	now := time.Now()
	idGeneratorMock := func() uuid.UUID {
		return id
	}
	createdAtGeneratorMock := func() time.Time {
		return now
	}
	todoDone1 := false
	todoDone2 := true
	todoDone3 := false
	wantedTodos := []model.Todo{
		{Id: uuid.New(), Title: "title1", Description: "description1", Done: &todoDone1, CreatedAt: time.Now()},
		{Id: uuid.New(), Title: "title2", Description: "description2", Done: &todoDone2, CreatedAt: time.Now()},
		{Id: uuid.New(), Title: "title3", Description: "description3", Done: &todoDone3, CreatedAt: time.Now()},
	}
	todoRepositoryImpl := TodoRepositoryImpl{DBPool: dbPoolMock,
		IDGenerator: idGeneratorMock, CreatedAtGenerator: createdAtGeneratorMock}
	gomock.InOrder(
		dbPoolMock.EXPECT().Query(gomock.Any(), queryAllTodos).Return(dbRowsMock, nil),
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
		dbRowsMock.EXPECT().Scan(gomock.Any()).SetArg(0, wantedTodos[2].CreatedAt).Return(anError),
	)
	todos, err := todoRepositoryImpl.GetAll()
	assert.Nil(t, todos)
	assert.Error(t, err, anError)
}

func TestGetAllWhenErrCallReturnAnError(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	dbRowsMock := common.NewMockDBRows(mockCtrl)
	id := uuid.New()
	now := time.Now()
	idGeneratorMock := func() uuid.UUID {
		return id
	}
	createdAtGeneratorMock := func() time.Time {
		return now
	}
	todoDone1 := false
	todoDone2 := true
	todoDone3 := false
	wantedTodos := []model.Todo{
		{Id: uuid.New(), Title: "title1", Description: "description1", Done: &todoDone1, CreatedAt: time.Now()},
		{Id: uuid.New(), Title: "title2", Description: "description2", Done: &todoDone2, CreatedAt: time.Now()},
		{Id: uuid.New(), Title: "title3", Description: "description3", Done: &todoDone3, CreatedAt: time.Now()},
	}
	todoRepositoryImpl := TodoRepositoryImpl{DBPool: dbPoolMock,
		IDGenerator: idGeneratorMock, CreatedAtGenerator: createdAtGeneratorMock}
	gomock.InOrder(
		dbPoolMock.EXPECT().Query(gomock.Any(), queryAllTodos).Return(dbRowsMock, anError),
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
		dbRowsMock.EXPECT().Err().Return(anError),
	)
	todos, err := todoRepositoryImpl.GetAll()
	assert.Nil(t, todos)
	assert.Error(t, err, anError)
}

func TestGetAllWhenIDGeneratorOrCreatedAtGeneratorIsNil(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	dbRowsMock := common.NewMockDBRows(mockCtrl)
	todoRepositoryImpl := TodoRepositoryImpl{DBPool: dbPoolMock,
		IDGenerator: nil, CreatedAtGenerator: nil}
	gomock.InOrder(
		dbPoolMock.EXPECT().Query(gomock.Any(), queryAllTodos).Return(dbRowsMock, nil),
		dbRowsMock.EXPECT().Next().Return(false),
		dbRowsMock.EXPECT().Err().Return(nil),
	)
	todos, err := todoRepositoryImpl.GetAll()
	assert.Equal(t, []model.Todo{}, todos)
	assert.NoError(t, err)
}

func TestGetAllWhenDBPoolIsNil(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	dbPoolMock := common.NewMockDBPool(mockCtrl)
	dbRowsMock := common.NewMockDBRows(mockCtrl)
	todoRepositoryImpl := TodoRepositoryImpl{DBPool: nil,
		IDGenerator: nil, CreatedAtGenerator: nil}
		dbPoolMock.EXPECT().Query(gomock.Any(), gomock.Any()).Times(0)
		dbRowsMock.EXPECT().Next().Times(0)
		dbRowsMock.EXPECT().Scan(gomock.Any()).Times(0)
		dbRowsMock.EXPECT().Err().Times(0)
	todos, err := todoRepositoryImpl.GetAll()
	assert.Nil(t, todos)
	assert.Equal(t, err, ErrTodoRepositoryInitialization)
}