package loggerTool

import (
	"fmt"
	"path"
	"strings"
	"time"
	"tools/appName"
	"tools/configTool"

	"github.com/fatih/color"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

func init() {
	// 配置输出格式
	logrus.SetFormatter(&MyFormatter{true})
	switch strings.ToLower(configTool.Log.Level) {
	case "panic", "p":
		logrus.SetLevel(logrus.PanicLevel)
	case "fatal", "f":
		logrus.SetLevel(logrus.FatalLevel)
	case "error", "err", "e":
		logrus.SetLevel(logrus.ErrorLevel)
	case "warning", "warn", "w":
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
	logPath := path.Join(configTool.Log.Dir, appName.Get())
	writer, err := rotatelogs.New(logPath+"_%Y%m%d.log",
		rotatelogs.WithLinkName(logPath+"_link"),
		rotatelogs.WithMaxAge(time.Duration(configTool.Log.MaxAge*24)*time.Hour),
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
	}, &MyFormatter{false})
	logrus.AddHook(hook)
}

type MyFormatter struct {
	isColor bool
}

// 根据不同log类型使用不同的输出样式
func (f *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	msg := entry.Message
	switch entry.Level {
	case logrus.TraceLevel:
		msg = fmt.Sprintln(timeText(entry), f.levelText(entry), ":", msg)
	case logrus.DebugLevel:
		msg = fmt.Sprintln(timeText(entry), f.levelText(entry), fileText(entry), ":", msg)
	case logrus.InfoLevel, logrus.WarnLevel:
		msg = fmt.Sprintln(timeText(entry), f.levelText(entry), functionText(entry), ":", msg)
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		msg = fmt.Sprintln(timeText(entry), f.levelText(entry), fileText(entry), functionText(entry), ":", msg)
	}
	return []byte(msg), nil
}

// 日志时间
func timeText(entry *logrus.Entry) string {
	return entry.Time.Format("2006-01-02 15:04:05")
}

// 日志等级
func (f *MyFormatter) levelText(entry *logrus.Entry) string {
	c := color.New()
	if f.isColor && !color.NoColor {
		c.EnableColor()
	} else {
		c.DisableColor()
	}
	switch entry.Level {
	case logrus.InfoLevel:
		return c.Add(color.FgBlue).Sprint(strings.ToUpper(entry.Level.String()))
	case logrus.DebugLevel:
		return c.Add(color.FgGreen).Sprint(strings.ToUpper(entry.Level.String()))
	case logrus.TraceLevel:
		return c.Add(color.FgCyan).Sprint(strings.ToUpper(entry.Level.String()))
	case logrus.WarnLevel:
		return c.Add(color.FgYellow).Sprint(strings.ToUpper(entry.Level.String()))
	case logrus.ErrorLevel, logrus.FatalLevel, logrus.PanicLevel:
		return c.Add(color.FgRed).Sprint(strings.ToUpper(entry.Level.String()))
	default:
		return strings.ToUpper(entry.Level.String())
	}
}

// 日志文件来源
func fileText(entry *logrus.Entry) string {
	return fmt.Sprintf("%s:%v", entry.Caller.File, entry.Caller.Line)
}

// 日志方法来源
func functionText(entry *logrus.Entry) string {
	functions := strings.Split(entry.Caller.Function, "/")
	return functions[len(functions)-1]
}
