package gormTool

import (
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
	"tools/configTool"

	"github.com/robfig/cron/v3"
)

type log struct {
	IsGorm bool   `default:"true"`
	Dir    string `default:"log"`
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
	if Log.IsGorm {
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
	change()
	c := cron.New()
	c.AddFunc("@daily", change)
	c.Start()
}

func change() {
	os.Mkdir("log", 0777)
	file, _ := os.OpenFile(filepath.Join(Log.Dir, time.Now().Format("gorm_20060102.log")), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	mw.change(file)
}
