package flagTool

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

var h = pflag.BoolP("help", "h", false, "帮助")
var v = pflag.BoolP("version", "v", false, "版本")

var version string
var goVersion string
var buildTime string

func HelpAndVersion() {
	if *h {
		pflag.PrintDefaults()
		os.Exit(0)
	}
	if *v {
		fmt.Println("Version:\t", version)
		fmt.Println("Go Version:\t", goVersion)
		fmt.Println("Build Time:\t", buildTime)
		os.Exit(0)
	}
}
