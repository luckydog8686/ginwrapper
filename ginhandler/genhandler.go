package ginhandler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/luckydog8686/errors"
	"github.com/luckydog8686/logs"
	"net/http"
	"reflect"
	"runtime"
	"strings"
)

func Default(c *gin.Context)  {
	c.JSON(http.StatusOK,gin.H{"error":"unsupportted api","data":""})
}

func Generate(f interface{})(map[string]gin.HandlerFunc,error)  {
	ftype := reflect.TypeOf(f)
	logs.Info(ftype.Kind())
	if ftype.Kind()== reflect.Func{
		return generateByFunc(f)
	}
	if ftype.Kind()==reflect.Struct{
		return generateByStruct(f)
	}
	if ftype.Kind()==reflect.Ptr && ftype.Elem().Kind()==reflect.Struct{
		logs.Info("generateByStructPtr")
		return generateByStructPtr(f)
	}

	return nil,nil
}

func  generateByFunc(f interface{})(map[string]gin.HandlerFunc,error){
	logs.Info("generateByFunc")
	ret := make(map[string]gin.HandlerFunc)
	ftype := reflect.TypeOf(f)
	numIn := ftype.NumIn()
	numOut := ftype.NumOut()
	paramsType:= make([]reflect.Type,0,numIn)
	for i:=0;i<numIn;i++{
		t := ftype.In(i)
		paramsType = append(paramsType,t)
	}
	outNames := make([]string,0,numOut)
	for i:=0;i<numOut;i++{
		outNames = append(outNames,ftype.Out(i).Name())
	}
	methodName := GetFunctionName(f)
	logs.Info("==",methodName,"==")
	ret[methodName]=func(context *gin.Context) {
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
					"data":nil,
				})
				return
			}
		}
		rst := call(f,params...)
		logs.Info(rst[0].Interface())
		context.JSON(http.StatusOK,gin.H{
			"data":rst[0].Interface(),
			"error":rst[1].Interface(),
		})
	}
	return ret,nil
}

func generateByStruct(s interface{})(map[string]gin.HandlerFunc,error)  {
	//ret := make(map[string]gin.HandlerFunc)

	val := reflect.ValueOf(s)
	logs.Info("generateByStruct::",val.NumMethod())

	for i :=0;i<val.NumMethod();i++{
		method := val.Type().Method(i)
		logs.Info(method.Type.NumIn())
		logs.Info(method.Type.In(0))
		logs.Info(method.Type.NumOut())
		logs.Info(method.Type.Out(0))
		logs.Info(method.Name)

	}
	return nil,nil
}

func generateByStructPtr(s interface{})(map[string]gin.HandlerFunc,error)  {
	vtype := reflect.TypeOf(s)
	structName := vtype.Elem().Name()
	val := reflect.ValueOf(s)
	if val.NumMethod() == 0{
		return nil,errors.New(fmt.Sprintf("no method for struct %s",structName))
	}
	logs.Info("generateByStruct::",val.NumMethod())
	ret := make(map[string]gin.HandlerFunc)

	for i :=0;i<val.NumMethod();i++{
		j := i
		method := val.Type().Method(j)
		mapKey := strings.ToLower(fmt.Sprintf("%s/%s",structName,method.Name))
		ret[mapKey]= func(c *gin.Context) {
			logs.Info(mapKey)
			numIn := method.Type.NumIn()
			var params []reflect.Value
			if numIn >1 {
				var s  reflect.Value
				if method.Type.In(1).Kind() == reflect.Ptr{
					s = reflect.New(method.Type.In(1).Elem())
				}else{
					s = reflect.New(method.Type.In(1))
				}

				s2 := s.Interface()
				if err:=c.Bind(s2);err!= nil{
					logs.Info(err)
					c.JSON(http.StatusOK,gin.H{
						"error":err,
						"data":nil,
					})
					return
				}
				if method.Type.In(1).Kind() == reflect.Ptr{
					params = append(params,s)
				}else{
					params = append(params,reflect.Indirect(s))
				}
			}
			rtns :=val.Method(j).Call(params)

			c.JSON(http.StatusOK,gin.H{
				"data":rtns[0].Interface(),
				"error":rtns[1].Interface(),
			})
		}
	}
	return ret,nil
}

func GetStructName(name string)string  {
	var seps []rune = []rune{'.','/'}
	fields := strings.FieldsFunc(name, func(sep rune) bool {
		for _, s := range seps {
			if sep == s {
				return true
			}
		}
		return false
	})

	// fmt.Println(fields)

	if size := len(fields); size > 0 {
		return strings.ToLower(fields[size-1])+"/"
	}
	return ""
}

func GetFunctionName(i interface{}) string {
	// 获取函数名称
	var seps []rune = []rune{'.','/'}
	fn := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()

	// 用 seps 进行分割
	fields := strings.FieldsFunc(fn, func(sep rune) bool {
		for _, s := range seps {
			if sep == s {
				return true
			}
		}
		return false
	})

	// fmt.Println(fields)

	if size := len(fields); size > 0 {
		return strings.ToLower(fields[size-1])
	}
	return ""
}

func  call(fun interface{},params ...interface{}) []reflect.Value {
	f := reflect.ValueOf(fun)
	in := make([]reflect.Value,len(params))
	for k, param := range params {
		in[k] = reflect.ValueOf(param)
	}
	return f.Call(in)
}