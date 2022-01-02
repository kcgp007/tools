package loggerInit

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
	"tools/configTool"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type log struct {
	Level string `default:"info"`
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
	w.Writer = io.MultiWriter(os.Stdout, file)
	w.File, file = file, w.File
	file.Close()
}

var mw = &myWriter{io.MultiWriter(os.Stdout), nil, new(sync.RWMutex)}

func init() {
	configTool.Add(&Log)
	enc := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())
	ws := zapcore.AddSync(mw)
	var enab zapcore.LevelEnabler
	switch strings.ToLower(Log.Level) {
	case "debug", "d":
		enab = zap.DebugLevel
	case "info", "i":
		enab = zap.InfoLevel
	case "warning", "warn", "w":
		enab = zap.WarnLevel
	case "error", "err", "e":
		enab = zap.ErrorLevel
	case "panic", "p":
		enab = zap.PanicLevel
	case "fatal", "f":
		enab = zap.FatalLevel
	default:
		enab = zap.InfoLevel
	}
	logger := zap.New(zapcore.NewCore(enc, ws, enab), zap.AddCaller())
	zap.ReplaceGlobals(logger)
	change()
	c := cron.New()
	c.AddFunc("@daily", change)
	c.Start()
}

func change() {
	os.Mkdir("log", os.ModePerm)
	file, _ := os.OpenFile(filepath.Join(Log.Dir, time.Now().Format("20060102.log")), os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	mw.change(file)
}
