package ginTool

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func Test_gin(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(gin.Logger(), ErrorHandler, ExceptionHandler)
	pprof.Register(router)
	router.GET("/ping", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})
	router.GET("/error", func(c *gin.Context) {
		c.Error(fmt.Errorf("test error"))
	})
	router.GET("/exception", func(c *gin.Context) {
		s := ""
		c.String(http.StatusOK, s[:5])
	})
	t.Log(router.Run(":8080"))
}
