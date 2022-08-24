package model

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

var validatorr *validator.Validate = validator.New()

type Todo struct {
	Id          uuid.UUID `json:"id"`
	Title       string    `json:"title" binding:"required" validate:"required"`
	Description string    `json:"description" binding:"required" validate:"required"`
	Done        *bool     `json:"done" binding:"required"`
	CreatedAt   time.Time `json:"createdAt"`
}

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
