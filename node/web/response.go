package web

import (
	"github.com/gin-gonic/gin"
	"strings"
)

type JsonObject map[string]interface{}

func respJSON(v interface{}) gin.H {
	return gin.H{
		"code": 0,
		"data": v,
	}
}
func respErrorCode(code int, c *gin.Context) gin.H {
	l := c.GetHeader("Lang")
	errSplit := strings.Split(ErrMap[code], ":")
	var e string
	switch l {
	case "cn":
		e = errSplit[1]
	default:
		e = errSplit[0]
	}
	return gin.H{
		"code": -1,
		"err":  code,
		"msg":  e,
	}
}
