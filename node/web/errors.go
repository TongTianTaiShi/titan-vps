package web

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"strings"
)

const (
	NotFound = iota + 1000
	InvalidParams
	UserNotFound

	Unknown     = -1
	GenericCode = 1
)

// ErrMap some errors from titan
var ErrMap = map[int]string{
	Unknown:       "unknown error:未知错误",
	NotFound:      "not found:信息未找到",
	InvalidParams: "invalid params:参数有误",
	UserNotFound:  "user not found:用户不存在",
}

type GenericError struct {
	Code int
	Err  error
}

func (e GenericError) Error() string {
	return e.Err.Error()
}

func NewErrorCode(Code int, c *gin.Context) GenericError {
	l := c.GetHeader("Lang")
	errSplit := strings.Split(ErrMap[Code], ":")
	var e string
	switch l {
	case "cn":
		e = errSplit[1]
	default:
		e = errSplit[0]
	}
	return GenericError{Code: Code, Err: errors.New(e)}

}
