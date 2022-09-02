package common

type AppError struct {
	Error   string `json:"error"`
	Message string `json:"message,omitempty"`
}

func SendBackAnAppError(ctx WebContext, log Logger, handlerError error, messege string, code int) {
	errObj := AppError{Error: handlerError.Error(), Message: messege}
	log.Printf("%v\n", errObj)
	ctx.JSON(code, errObj)
}
