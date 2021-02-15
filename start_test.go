package ginwrapper

import (
	"github.com/gin-gonic/gin"
	"github.com/luckydog8686/ginwrapper/ginhandler"
	"testing"
)
type SS struct {
	Name string
}

func Ping(s *SS) gin.H{
	ret := make(map[string]interface{})
	ret["error"]=nil
	ret["data"]=s
	return ret
}
func TestStart(t *testing.T) {
	var tmp map[string]gin.HandlerFunc = make(map[string]gin.HandlerFunc)
	var err error
	if tmp["ping"],err =ginhandler.Generate(Ping);err != nil{
		t.Fatal(err)
	}
	Start(tmp,"127.0.0.1:80")
}
