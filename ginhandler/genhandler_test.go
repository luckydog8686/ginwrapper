package ginhandler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/luckydog8686/errors"
	"testing"
)

func TestGenerate(t *testing.T) {
	f,err := Generate(Ping)
	if err != nil {
		t.Fatal(err)
	}
	router := gin.Default()
	for k,v := range f{
		router.POST(fmt.Sprintf("%s",k),v)
	}

	router.Run("127.0.0.1:80")
}








type SS struct {
	Name string
}

func (s *SS)Hello(str string) (string,error)   {
	return "fuck the world",errors.New("fuck")
}


func Ping(s *SS) (*SS,error){
	return s,errors.New("fuck")

}