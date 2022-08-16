package model

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	Id uuid.UUID `json:"id"`
	Title string `json:"title" binding:"required"`
	Description string `json:"description" binding:"required"`
	Done *bool `json:"done" binding:"required"`
	CreatedAt time.Time `json:"createdAt"`
}