package ginhandler

import (
	"github.com/gin-gonic/gin"
	"testing"
)

func TestGenerate(t *testing.T) {
	f,err := Generate(Ping)
	if err != nil {
		t.Fatal(err)
	}
	router := gin.Default()
	router.POST("/ping",f)
	router.Run("127.0.0.1:80")
}

type SS struct {
	Name string
}

func Ping(s *SS) gin.H{
	ret := make(map[string]interface{})
	ret["error"]=nil
	ret["data"]=s
	return ret
}