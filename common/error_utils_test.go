package common

import (
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestDisplayAppError(t *testing.T) {
	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	log := LoggerMock{}
	handlerErr := errors.New("handlerErr1")
	message := "message1"
	code := float64(404)
	DisplayAppError(ctx, &log, handlerErr, message, code)
	if w.Code != int(code) {
		t.Errorf("Expects %d, but got %d", int(code), w.Code)
	}

	var got gin.H
	err := json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}
	want := gin.H{"error": handlerErr.Error(), "message": message}
	assert.Equal(t, want, got)
	assert.Equal(t, 1, log.Called)
}