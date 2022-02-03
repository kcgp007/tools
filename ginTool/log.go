package ginTool

import (
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kcgp007/tools/configTool"
	"github.com/robfig/cron/v3"
)

type log struct {
	IsGin bool   `default:"true"`
	Dir   string `default:"log"`
}

var Log log

type myWriter struct {
	io.Writer
	*os.File
	*sync.RWMutex
}

func (w *myWriter) Write(p []byte) (n int, err error) {
	w.RLock()
	defer w.RUnlock()
	return w.Writer.Write(p)
}

func (w myWriter) change(file *os.File) {
	w.Lock()
	defer w.Unlock()
	if Log.IsGin {
		w.Writer = io.MultiWriter(os.Stdout, file)
	} else {
		w.Writer = io.MultiWriter(file)
	}
	w.File, file = file, w.File
	file.Close()
}

var mw = &myWriter{io.MultiWriter(os.Stdout), nil, new(sync.RWMutex)}

func init() {
	configTool.Add(&Log)
	gin.DisableConsoleColor()
	gin.DefaultWriter = mw
	change()
	c := cron.New()
	c.AddFunc("@daily", change)
	c.Start()
}

func change() {
	os.Mkdir(Log.Dir, os.ModePerm)
	file, _ := os.OpenFile(filepath.Join(Log.Dir, time.Now().Format("gin_20060102.log")), os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	mw.change(file)
}
