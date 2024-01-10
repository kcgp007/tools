# tools

[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/kcgp007/tools/go.yml?logo=github)](https://github.com/kcgp007/tools/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/kcgp007/tools)](https://goreportcard.com/report/github.com/kcgp007/tools)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/kcgp007/tools)
![GitHub tag (latest by date)](https://img.shields.io/github/v/tag/kcgp007/tools)
![GitHub](https://img.shields.io/github/license/kcgp007/tools)
[![Go Dev](https://img.shields.io/badge/go.dev-reference-brightgreen?logo=go)](https://pkg.go.dev/github.com/kcgp007/tools)
![GitHub last commit](https://img.shields.io/github/last-commit/kcgp007/tools)
![GitHub commit activity](https://img.shields.io/github/commit-activity/m/kcgp007/tools)

tools 是一个 Go 编写的便携工具库，包含一些初始化功能和默认配置功能，同时将日志、配置、分析工具集成到了一起。
- 随机种子的初始化
- [zap](https://github.com/uber-go/zap) 日志库的初始化
- [viper](https://github.com/spf13/viper) 默认配置
- [pflag](https://github.com/spf13/pflag) 默认配置
- [gin](https://github.com/gin-gonic/gin) 默认配置
- [gorm](https://gorm.io/) 默认配置

# 初始化

```go
import (
	_ "github.com/kcgp007/tools/randInit"
	_ "github.com/kcgp007/tools/loggerInit"
	_ "github.com/kcgp007/tools/gormTool"
)
```

# 默认配置

## viper
```go
type dataConfig struct {
	S string
	I int
	F float64
}

var data = dataConfig{
	S: "aaa",
	I: 10,
	F: 10.5,
}

// 添加配置
configTool.Add(&data)
// 添加配置并自定义配置名
configTool.AddWithKey("data", &data)
```

## pflag
```go
// 添加方法
flagTool.HelpAndVersion()
```
```Makefile
// 编译时添加参数
go build -trimpath -ldflags "-s -w -X 'github.com/kcgp007/tools/flagTool.version=$(GIT_VERSION)' -X 'github.com/kcgp007/tools/flagTool.goVersion=$(GO_VERSION)' -X 'github.com/kcgp007/tools/flagTool.buildTime=$(BUILD_TIME)'"
```

## gin
```go
gin.SetMode(gin.TestMode)
router := gin.New()
// 添加错误处理和拦截
router.Use(gin.Logger(), ErrorHandler, ExceptionHandler)
// 添加分析工具
pprof.Register(router)
router.GET("/ping", func(c *gin.Context) {
	c.Status(http.StatusOK)
})
router.GET("/error", func(c *gin.Context) {
	c.Error(fmt.Errorf("test error"))
})
router.GET("/exception", func(c *gin.Context) {
	s := ""
	c.String(http.StatusOK, s[:5])
})
router.Run(":8080")
```

# 感谢 JetBrains 的支持

![JetBrains Logo (Main) logo](https://resources.jetbrains.com/storage/products/company/brand/logos/jb_beam.svg)

![GoLand logo](https://resources.jetbrains.com/storage/products/company/brand/logos/GoLand_icon.svg)
