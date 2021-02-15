package ginwrapper

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

func Start(handlers map[string]gin.HandlerFunc,listen string)  {
	r := gin.Default()
	for k,v := range handlers{
		r.POST(fmt.Sprintf("/%s",k),v)
	}
	r.Run(listen)
}
