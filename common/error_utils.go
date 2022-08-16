package common

import (
	"github.com/gin-gonic/gin"
)

type AppError struct {
	Error string `json:"error"`
	Message string `json:"message,omitempty"`
}

func DisplayAppError(ctx *gin.Context, log Logger, handlerError error, messege string, code float64) {
	errObj := AppError{Error: handlerError.Error(), Message: messege}
	log.Printf("%v\n", errObj)
	ctx.JSON(int(code), errObj)
}