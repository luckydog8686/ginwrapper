package cmd

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/luckydog8686/ginlogger"
)

func Run(listenAddr string,handlers map[string]gin.HandlerFunc)  {
	router := gin.New()
	router.Use(ginlogger.Logger(), gin.Recovery())
	for k,v := range handlers {
		router.POST(fmt.Sprintf("/%s",k),v)
	}
	router.Run(listenAddr)
}