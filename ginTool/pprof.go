package ginTool

import (
	"net/http"
	"net/http/pprof"

	"github.com/gin-gonic/gin"
)

func PprofRouter(router *gin.Engine) {
	group := router.Group("/debug/pprof")
	group.GET("/", pprofHandler(pprof.Index))
	group.GET("/cmdline", pprofHandler(pprof.Cmdline))
	group.GET("/profile", pprofHandler(pprof.Profile))
	group.POST("/symbol", pprofHandler(pprof.Symbol))
	group.GET("/trace", pprofHandler(pprof.Trace))
	group.GET("/allocs", pprofHandler(pprof.Handler("allocs").ServeHTTP))
	group.GET("/block", pprofHandler(pprof.Handler("block").ServeHTTP))
	group.GET("/goroutine", pprofHandler(pprof.Handler("goroutine").ServeHTTP))
	group.GET("/heap", pprofHandler(pprof.Handler("heap").ServeHTTP))
	group.GET("/mutex", pprofHandler(pprof.Handler("mutex").ServeHTTP))
	group.GET("/threadcreate", pprofHandler(pprof.Handler("threadcreate").ServeHTTP))
}

func pprofHandler(handlerFunc http.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		handlerFunc.ServeHTTP(c.Writer, c.Request)
	}
}
