package appName

import (
	"os"
	"os/exec"
	"path/filepath"
)

func Get() string {
	appName, _ := exec.LookPath(os.Args[0])
	ext := filepath.Ext(appName)
	appName = filepath.Base(appName)
	appName = appName[:len(appName)-len(ext)]
	return appName
}
