package model

import (
	"time"

	"github.com/go-playground/validator/v10"
)

var validatorr *validator.Validate = validator.New()

type Todo struct {
	Id          string    `json:"id" bindging:"required,uuid4" validate:"required,uuid4"`
	Title       string    `json:"title" binding:"required" validate:"required"`
	Description string    `json:"description" binding:"required" validate:"required"`
	Done        *bool     `json:"done" binding:"required" validate:"required"`
	CreatedAt   time.Time `json:"createdAt"`
}

//TODO Review the created Todos which doesn't have Id field
func IsValid(obj interface{}) (ok bool) {
	if obj == nil {
		return false
	}
	err := validatorr.Struct(obj)
	if err != nil {
		return false
	}
	return true
}
