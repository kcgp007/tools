package tools

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/fatih/color"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

func init() {
	myFormatter := new(MyFormatter)
	// 配置输出格式
	logrus.SetFormatter(myFormatter)
	switch strings.ToLower(Log.Level) {
	case "panic", "p":
		logrus.SetLevel(logrus.PanicLevel)
	case "fatal", "f":
		logrus.SetLevel(logrus.FatalLevel)
	case "error", "e", "err":
		logrus.SetLevel(logrus.ErrorLevel)
	case "warn", "warning", "w":
		logrus.SetLevel(logrus.WarnLevel)
	case "info", "i":
		logrus.SetLevel(logrus.InfoLevel)
	case "debug", "d":
		logrus.SetLevel(logrus.DebugLevel)
	case "trace", "t":
		logrus.SetLevel(logrus.TraceLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}
	logrus.SetReportCaller(true)
	// 获取配置文件中的日志文件路径
	appName, _ := exec.LookPath(os.Args[0])
	ext := filepath.Ext(appName)
	appName = filepath.Base(appName)
	appName = appName[:len(appName)-len(ext)]
	logPath := path.Join(Log.Dir, appName)
	writer, err := rotatelogs.New(logPath+"_%Y%m%d.log",
		rotatelogs.WithLinkName(logPath+"_link"),
		rotatelogs.WithMaxAge(time.Duration(Log.MaxAge*24)*time.Hour),
		rotatelogs.WithRotationTime(time.Duration(24)*time.Hour))
	if err != nil {
		logrus.Panic("init logger error:", err)
	}
	hook := lfshook.NewHook(lfshook.WriterMap{
		logrus.TraceLevel: writer,
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
	}, myFormatter)
	logrus.AddHook(hook)
}

type MyFormatter struct{}

// 根据不同log类型使用不同的输出样式
func (_ *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	msg := entry.Message
	switch entry.Level {
	case logrus.TraceLevel:
		msg = fmt.Sprintln(timeText(entry), levelText(entry), ":", msg)
	case logrus.DebugLevel:
		msg = fmt.Sprintln(timeText(entry), levelText(entry), fileText(entry), ":", msg)
	case logrus.InfoLevel, logrus.WarnLevel:
		msg = fmt.Sprintln(timeText(entry), levelText(entry), functionText(entry), ":", msg)
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		msg = fmt.Sprintln(timeText(entry), levelText(entry), fileText(entry), functionText(entry), ":", msg)
	}
	return []byte(msg), nil
}

// 日志时间
func timeText(entry *logrus.Entry) string {
	return entry.Time.Format("2006-01-02 15:04:05")
}

// 日志等级
func levelText(entry *logrus.Entry) string {
	switch entry.Level {
	case logrus.InfoLevel:
		return color.New(color.FgBlue).Sprint(strings.ToUpper(entry.Level.String()))
	case logrus.DebugLevel, logrus.TraceLevel:
		return color.New(color.FgGreen).Sprint(strings.ToUpper(entry.Level.String()))
	case logrus.WarnLevel:
		return color.New(color.FgYellow).Sprint(strings.ToUpper(entry.Level.String()))
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return color.New(color.FgRed).Sprint(strings.ToUpper(entry.Level.String()))
	default:
		return strings.ToUpper(entry.Level.String())
	}
}

// 日志来源
func fileText(entry *logrus.Entry) string {
	return fmt.Sprintf("%s:%v", entry.Caller.File, entry.Caller.Line)
}

// 日志内容
func functionText(entry *logrus.Entry) string {
	return entry.Caller.Function
}
