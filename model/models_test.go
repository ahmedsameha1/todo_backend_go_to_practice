package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsValidWhenTodoIsNil(t *testing.T) {
	var todo interface{} = nil
	ok := IsValid(todo)
	assert.False(t, ok)
}

func TestIsValidWhenTodoTitleIsEmpty(t *testing.T) {
	todo := Todo{}
	ok := IsValid(todo)
	assert.False(t, ok)
}

func TestIsValidWhenTodoTitleIsEmpty2(t *testing.T) {
	todo := Todo{Title: ""}
	ok := IsValid(todo)
	assert.False(t, ok)
}

func TestIsValidWhenTodoDescriptionIsEmpty(t *testing.T) {
	todo := Todo{}
	ok := IsValid(todo)
	assert.False(t, ok)
}

func TestIsValidWhenTodoDescriptionIsEmpty2(t *testing.T) {
	todo := Todo{Description: ""}
	ok := IsValid(todo)
	assert.False(t, ok)
}

func TestIsValidWhenTodoIsValid(t *testing.T) {
	todo := Todo{Description: "description", Title: "title"}
	ok := IsValid(todo)
	assert.True(t, ok)
}
