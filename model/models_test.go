package model

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestIsValidWhenTodoIsNil(t *testing.T) {
	var todo interface{} = nil
	ok := IsValid(todo)
	assert.False(t, ok)
}

func TestIsValidWhenTodoTitleIsEmpty(t *testing.T) {
	todoDone := false
	todo := Todo{Description: "description", Done: &todoDone}
	ok := IsValid(todo)
	assert.False(t, ok)
}

func TestIsValidWhenTodoTitleIsEmpty2(t *testing.T) {
	todoDone := false
	todo := Todo{Title: "", Description: "description", Done: &todoDone}
	ok := IsValid(todo)
	assert.False(t, ok)
}

func TestIsValidWhenTodoDescriptionIsEmpty(t *testing.T) {
	todoDone := false
	todo := Todo{Title: "title", Description: "", Done: &todoDone}
	ok := IsValid(todo)
	assert.False(t, ok)
}

func TestIsValidWhenTodoDescriptionIsEmpty2(t *testing.T) {
	todoDone := false
	todo := Todo{Title: "title", Done: &todoDone}
	ok := IsValid(todo)
	assert.False(t, ok)
}

func TestIsValidWhenTodoDoneIsNil(t *testing.T) {
	todo := Todo{Title: "titlde", Description: "description", Done: nil}
	ok := IsValid(todo)
	assert.False(t, ok)
}

func TestIsValidWhenTodoIsValid(t *testing.T) {
	todoDone := false
	todo := Todo{Id: uuid.New().String(), Description: "description", Title: "title", Done: &todoDone}
	ok := IsValid(todo)
	assert.True(t, ok)
}

func TestIsValid(t *testing.T) {
	t.Run("When there is no id", func(t *testing.T) {
		todoDone := false
		todo := Todo{Description: "description", Title: "title", Done: &todoDone}
		assert.False(t, IsValid(todo))
	})
}
