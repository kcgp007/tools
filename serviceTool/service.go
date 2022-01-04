package serviceTool

import (
	"fmt"
	"log"

	"github.com/kardianos/service"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
)

var start = pflag.Bool("start", false, "开始服务")
var stop = pflag.Bool("stop", false, "停止服务")
var restart = pflag.Bool("restart", false, "重启服务")
var install = pflag.Bool("install", false, "注册服务")
var uninstall = pflag.Bool("uninstall", false, "卸载服务")
var name = pflag.Bool("name", false, "服务名称")
var status = pflag.Bool("status", false, "服务状态")

func ServiceSet(s service.Service, err error) {
	if err != nil {
		zap.S().Panic(err)
	}
	switch {
	case *start:
		err = s.Start()
		if err != nil {
			zap.S().Error(err)
		}
	case *stop:
		err = s.Stop()
		if err != nil {
			log.Panicln(err)
		}
	case *restart:
		err = s.Restart()
		if err != nil {
			log.Panicln(err)
		}
	case *install:
		err = s.Install()
		if err != nil {
			log.Panicln(err)
		}
	case *uninstall:
		err = s.Uninstall()
		if err != nil {
			log.Panicln(err)
		}
	case *name:
		fmt.Println(s.String())
	case *status:
		st, err := s.Status()
		if err != nil {
			log.Panicln(err)
		}
		switch st {
		case service.StatusUnknown:
			fmt.Println("服务状态未知")
		case service.StatusRunning:
			fmt.Println("服务运行中")
		case service.StatusStopped:
			fmt.Println("服务已停止")
		}
	default:
		err = s.Run()
		if err != nil {
			log.Panicln(err)
		}
	}
}
