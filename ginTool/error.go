package ginTool

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 错误处理
func ErrorHandler(c *gin.Context) {
	c.Next()
	if err := c.Errors.Last(); err != nil {
		switch err.Type {
		case gin.ErrorTypeBind:
			c.String(http.StatusBadRequest, err.Error())
		default:
			c.String(http.StatusInternalServerError, err.Error())
		}
	}
}

// 异常处理
func ExceptionHandler(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			logrus.Error(err)
			c.String(http.StatusInternalServerError, fmt.Sprint(err))
		}
	}()
	c.Next()
}
