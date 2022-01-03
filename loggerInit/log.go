package loggerInit

import (
	"io"
	"os"
	"path/filepath"
	"strings"
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

var (
	tempLogger *zap.Logger
	tempFile   *os.File
)

func init() {
	configTool.Add(&Log)
	change()
	c := cron.New()
	c.AddFunc("@daily", change)
	c.Start()
}

func change() {
	enc := zapcore.NewJSONEncoder(zap.NewDevelopmentEncoderConfig())
	os.Mkdir("log", os.ModePerm)
	file, _ := os.OpenFile(filepath.Join(Log.Dir, time.Now().Format("20060102.log")), os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	ws := zapcore.AddSync(io.MultiWriter(file))
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

	if tempLogger != nil {
		tempLogger.Sync()
	}
	if tempFile != nil {
		tempFile.Close()
	}
	tempLogger, tempFile = logger, file
}
