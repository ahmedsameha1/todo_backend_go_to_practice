package common

import (
	"github.com/gin-gonic/gin"
)

type appError struct {
	Error string `json:"error"`
	Message string `json:"message"`
	HttpStatus float64 `json:"Status"`
}

func DisplayAppError(ctx *gin.Context, log Logger, handlerError error, messege string, code float64) {
	errObj := appError{HttpStatus: code, Error: handlerError.Error(), Message: messege}
	log.Printf("appError %s\n", errObj.Error)
	ctx.JSON(int(code), errObj)
}