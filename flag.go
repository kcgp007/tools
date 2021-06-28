package tools

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

var help = pflag.BoolP("help", "h", false, "帮助")
var v = pflag.BoolP("version", "v", false, "版本")

// CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -trimpath -ldflags "-X main.version=v1.2.3"
var version string

// 帮助&版本
func init() {
	pflag.Parse()
	if *help {
		pflag.PrintDefaults()
		os.Exit(0)
	}
	if *v {
		fmt.Println(version)
		os.Exit(0)
	}
}
