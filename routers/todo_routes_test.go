package routers

import (
	"testing"
	"github.com/ahmedsameha1/todo_backend_go_to_practice/common"
)

func TestSetTodoRoutes(t *testing.T){
	var poster common.PosterMock = common.PosterMock{}
	SetTodoRoutes(&poster)
	if poster.Called != 6 {
		t.Errorf("Expects %d but got %d", 6, poster.Called )
	}
}