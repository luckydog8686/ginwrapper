package ginhandler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/luckydog8686/errors"
	"github.com/luckydog8686/logs"
	"net/http"
	"testing"
)

func TestGenerate(t *testing.T) {
	s := &SS{}
	f,err := Generate(s)
	if err != nil {
		t.Fatal(err)
	}
	router := gin.Default()
	for k,v := range f{
		router.POST(fmt.Sprintf("%s",k),v)
	}
	router.POST("fuck/you", func(context *gin.Context) {
		context.JSON(http.StatusOK,"fuckkk")
	})
	router.Run("127.0.0.1:80")
}








type SS struct {
	Name string
}

func (s *SS)Hello(str *SS) (string,error)   {
	logs.Info("=======",str.Name)
	return str.Name+"fuck the world",errors.New("fuck")
}
func (s *SS)Ping(str string) (string,error)   {
	logs.Info(str)
	return fmt.Sprintf("ping %s",str),nil
}


func Ping(s *SS) (*SS,error){
	return s,errors.New("fuck")

}