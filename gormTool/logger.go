package gormTool

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm/logger"
)

func (w *myWriter) LogMode(logger.LogLevel) logger.Interface {
	return w
}

func (w *myWriter) Info(_ context.Context, msg string, data ...interface{}) {
	w.WriteString(fmt.Sprintf("%v [Info] %v", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(msg, data...)))
}

func (w *myWriter) Warn(_ context.Context, msg string, data ...interface{}) {
	w.WriteString(fmt.Sprintf("%v [Warn] %v", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(msg, data...)))
}

func (w *myWriter) Error(_ context.Context, msg string, data ...interface{}) {
	w.WriteString(fmt.Sprintf("%v [Error] %v", time.Now().Format("2006-01-02 15:04:05"), fmt.Sprintf(msg, data...)))
}

func (w *myWriter) Trace(_ context.Context, begin time.Time, fc func() (string, int64), _ error) {
	sql, rows := fc()
	w.WriteString(fmt.Sprintf("%v:%v [%v] | %v", begin.Format("2006-01-02 15:04:05"), time.Since(begin), rows, sql))
}
