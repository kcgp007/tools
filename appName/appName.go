package appName

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

var appName string
var once sync.Once

func Get() string {
	once.Do(func() {
		appName, _ = exec.LookPath(os.Args[0])
		ext := filepath.Ext(appName)
		appName = filepath.Base(appName)
		appName = appName[:len(appName)-len(ext)]
	})
	return appName
}
