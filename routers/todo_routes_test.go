package routers

import (
	"testing"

)

func TestSetTodoRoutes(t *testing.T){
	var poster PosterMock = PosterMock{}
	SetTodoRoutes(&poster)
	if poster.called != 6 {
		t.Errorf("Expects %d but got %d", 6, poster.called )
	}
}