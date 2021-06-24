package tools

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

var help = pflag.BoolP("help", "h", false, "帮助")
var v = pflag.BoolP("version", "v", false, "版本")

// -ldflags "-X tools.version=v1.2.3"
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
