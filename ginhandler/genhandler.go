package ginhandler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/luckydog8686/logs"
	"net/http"
	"reflect"
)

func Default(c *gin.Context)  {
	c.JSON(http.StatusOK,gin.H{"error":"unsupportted api","data":""})
}

func Generate(f interface{})(gin.HandlerFunc,error)  {
	ftype := reflect.TypeOf(f)
	if ftype.Kind()!= reflect.Func{
		return Default,errors.New("Not a func method")
	}
	numIn := ftype.NumIn()
	numOut := ftype.NumOut()
	logs.Info(numIn)
	paramsType:= make([]reflect.Type,0,numIn)
	logs.Info(len(paramsType))
	for i:=0;i<numIn;i++{
		t := ftype.In(i)
		paramsType = append(paramsType,t)
	}
	logs.Info(len(paramsType))
	outNames := make([]string,0,numOut)
	for i:=0;i<numOut;i++{
		outNames = append(outNames,ftype.Out(i).Name())
	}

	return func(context *gin.Context) {
		params := make([]interface{},0,numIn)
		logs.Info(params)
		for j:=0;j<numIn;j++{
			if paramsType[j].Kind() != reflect.Ptr{
				ifc := reflect.New(paramsType[j]).Interface()
				params = append(params,ifc)
			}else{
				ifc := reflect.New(paramsType[j].Elem()).Interface()
				params = append(params,ifc)
			}
		}
		if numIn>0{
			err := context.Bind(params[0])
			if err !=nil{
				context.JSON(http.StatusOK,gin.H{
					"error":err,
					"data":"",
				})
				return
			}
		}
		rst := call(f,params...)

		context.JSON(http.StatusOK,rst[0].Interface())
	},nil
}

func  call(fun interface{},params ...interface{}) []reflect.Value {
	f := reflect.ValueOf(fun)
	in := make([]reflect.Value,len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	return f.Call(in)
}