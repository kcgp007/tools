package app

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

var name string
var once sync.Once

func Name() string {
	once.Do(func() {
		name, _ = exec.LookPath(os.Args[0])
		ext := filepath.Ext(name)
		name = filepath.Base(name)
		name = name[:len(name)-len(ext)]
	})
	return name
}
