package flagInit

import (
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

var v = pflag.BoolP("version", "v", false, "版本")

// CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc
// CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=x86_64-linux-musl-gcc
// CGO_ENABLED=1 GOOS=linux GOARCH=arm CC=arm-linux-gnueabihf-gcc
// go build -trimpath -ldflags "-s -w -X tools.version=v1.2.3"
var version string

// 帮助&版本
func init() {
	pflag.Parse()
	if *v {
		fmt.Println(version)
		os.Exit(0)
	}
}
