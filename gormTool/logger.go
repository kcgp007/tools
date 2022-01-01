package gormTool

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm/logger"
)

type MyLogger struct {
	myWriter
}

func (MyLogger *MyLogger) LogMode(logger.LogLevel) logger.Interface {
	return MyLogger
}

func (MyLogger *MyLogger) Info(_ context.Context, msg string, data ...interface{}) {
	MyLogger.WriteString(fmt.Sprintf("%v [Info] %v", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(msg, data)))
}

func (MyLogger *MyLogger) Warn(_ context.Context, msg string, data ...interface{}) {
	MyLogger.WriteString(fmt.Sprintf("%v [Warn] %v", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(msg, data)))
}

func (MyLogger *MyLogger) Error(_ context.Context, msg string, data ...interface{}) {
	MyLogger.WriteString(fmt.Sprintf("%v [Error] %v", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(msg, data)))
}

func (MyLogger *MyLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), _ error) {
	sql, rows := fc()
	MyLogger.WriteString(fmt.Sprintf("%v:%v [%v] | %v", begin.Format("2006-01-02 15:04:05"), time.Since(begin), rows, sql))
}
