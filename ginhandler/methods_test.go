package ginhandler

import (
	"github.com/gin-gonic/gin"
	"github.com/luckydog8686/logs"
	"testing"
)

func Funciuii1()  {
	logs.Info("Func1 print")
}


type FStct struct {
}

func (fs *FStct)Print(ctx *gin.Context)  {
	logs.Info("fs Print")
}

func TestExtractMethods(t *testing.T) {
	if _,err := ExtractMethods(new(FStct));err != nil {
		t.Fatal(err)
	}
}