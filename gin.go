package go_utils

import (
	"encoding/json"
	"io"
	"strings"

	"errors"
	"github.com/gin-gonic/gin"
	"reflect"

	"github.com/go-playground/validator/v10"
)

type baseController struct {
}

// 绑定模型获取验证错误的方法
func (c *baseController) GetEventError(ctx *gin.Context, errs validator.ValidationErrors) string {

	for _, e := range errs {
		ns := e.Namespace()
		i := strings.LastIndex(ns, ".") + 1
		if i > 0 {
			ss := ns[i:]
			ss = UnMarshal(ss)
			return ss + "不能为空"
		}
	}

	return "参数错误"
}

// 获取错误信息
func (c *baseController) GetError(st interface{}, errs validator.ValidationErrors) string {
	for _, e := range errs {
		field, _ := reflect.TypeOf(st).Elem().FieldByName(e.Field())
		error := field.Tag.Get("error")
		if error != "" {
			return error
		}
	}
	return "参数错误"
}


// 校验参数
func Check(requ interface{}, ctx *gin.Context) error {
	if err := ctx.ShouldBind(requ); err != nil {
		if err == io.EOF {
			return NewParamError(400, errors.New("请求参数不能为空"), err)
		}
		if e, ok := err.(validator.ValidationErrors); ok {
			eventError := GetError(requ, e)
			return NewParamError(400, errors.New(eventError), nil)
		}
		if _, ok := err.(*json.SyntaxError); ok {
			return NewParamError(400, errors.New("请求数据格式错误"), err)
		}
		return NewParamError(400, errors.New("未知错误"), err)
	}
	return nil
}


// 获取错误信息
func GetError(st interface{}, errs validator.ValidationErrors) string {
	for _, e := range errs {
		field, _ := reflect.TypeOf(st).Elem().FieldByName(e.Field())
		error := field.Tag.Get("error")
		if error != "" {
			return error
		}
	}
	return "参数错误"
}

