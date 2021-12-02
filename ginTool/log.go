package ginTool

import (
	"io"
	"os"
	"tools/loggerInit"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.DisableConsoleColor()
	gin.DefaultWriter = io.MultiWriter(os.Stdout, loggerInit.Writer)
}
