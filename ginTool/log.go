package ginTool

import (
	"io"
	"os"
	"tools/loggerTool"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.DisableConsoleColor()
	gin.DefaultWriter = io.MultiWriter(os.Stdout, loggerTool.Writer)
}
