package gormTool

import (
	"context"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm/logger"
)

type MyLogger struct{}

func (MyLogger *MyLogger) LogMode(logger.LogLevel) logger.Interface {
	return MyLogger
}

func (MyLogger *MyLogger) Info(_ context.Context, _ string, data ...interface{}) {
	logrus.Info(data)
}

func (MyLogger *MyLogger) Warn(_ context.Context, _ string, data ...interface{}) {
	logrus.Warn(data)
}

func (MyLogger *MyLogger) Error(_ context.Context, _ string, data ...interface{}) {
	logrus.Error(data)
}

func (MyLogger *MyLogger) Trace(_ context.Context, begin time.Time, fc func() (string, int64), _ error) {
	sql, rows := fc()
	logrus.Tracef("%v [%v] | %v", time.Since(begin), rows, sql)
}
