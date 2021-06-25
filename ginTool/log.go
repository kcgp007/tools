package ginTool

import (
	"net/http"
	"strconv"
	"time"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// 访问日志
func WebLog(c *gin.Context) {
	startTime := time.Now()
	c.Next()
	latencyTime := time.Since(startTime)
	logrus.Infof("[%v] %15v | %v %v %#v",
		c.ClientIP(),
		latencyTime,
		statusText(c.Writer.Status()),
		methodText(c.Request.Method),
		c.Request.RequestURI)
}

// HTTP 状态码
func statusText(statusCode int) string {
	switch {
	case statusCode >= http.StatusOK && statusCode < http.StatusMultipleChoices:
		return color.New(color.BgGreen).Sprint(statusCode)
	case statusCode >= http.StatusMultipleChoices && statusCode < http.StatusBadRequest:
		return color.New(color.BgBlue).Sprint(statusCode)
	case statusCode >= http.StatusBadRequest && statusCode < http.StatusInternalServerError:
		return color.New(color.BgYellow).Sprint(statusCode)
	case statusCode >= http.StatusInternalServerError:
		return color.New(color.BgRed).Sprint(statusCode)
	default:
		return strconv.Itoa(statusCode)
	}
}

// HTTP 请求方法
func methodText(method string) string {
	switch method {
	case http.MethodGet:
		return color.New(color.BgBlue).Sprint(method)
	case http.MethodPost:
		return color.New(color.BgCyan).Sprint(method)
	case http.MethodPut:
		return color.New(color.BgYellow).Sprint(method)
	case http.MethodDelete:
		return color.New(color.BgRed).Sprint(method)
	case http.MethodPatch:
		return color.New(color.BgGreen).Sprint(method)
	case http.MethodHead:
		return color.New(color.BgMagenta).Sprint(method)
	case http.MethodOptions:
		return color.New(color.BgWhite, color.FgBlack).Sprint(method)
	default:
		return method
	}
}
