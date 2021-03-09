package ginhandler

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/luckydog8686/logs"
	"path/filepath"
	"reflect"
	"runtime"
)

func ExtractMethods(params ...interface{}) (map[string]interface{},error) {
	funcMap := make(map[string]gin.HandlerFunc)
	for _,v := range params{
		vt:=reflect.TypeOf(v)
		if vt.Kind()==reflect.Struct{
			logs.Info("struct:",vt.NumMethod())
			for i:=0;i<vt.NumMethod();i++{
				mtd := vt.Method(i)
				mtdt:=reflect.TypeOf(mtd.Func.Elem().Interface())

				logs.Info(mtd.Name,"==========",mtd.Type,"===========",mtdt)
			}
			continue
		}
		if vt.Kind()==reflect.Ptr{
			logs.Info("ptr")
		}
		if vt.Kind()==reflect.Func{
			handler,err := Generate(v)

			if err != nil {
				logs.Fatal(err)
			}
			for k,v := range handler{
				funcMap[k]=v
			}
			continue
		}
		return nil,errors.New("Encounter unsupported kind, only struct and func are supported!")
	}
	return nil,nil
}

func GetFuncName(i interface{}) string {
	return string([]byte(filepath.Ext(runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()))[1:])
}