package model

import (
	"time"

	"github.com/google/uuid"
)

type Todo struct {
	Id uuid.UUID `json:"id"`
	Title string `json:"title"`
	Description string `json:"description"`
	Done bool `json:"done"`
	CreatedAt time.Time `json:"createdAt"`
}