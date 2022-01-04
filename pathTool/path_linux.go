package pathTool

import (
	"os"
	"path/filepath"
)

func SmartWd() string {
	wd, _ := os.Getwd()
	return wd
}

func SmartAbs(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(SmartWd(), path)
}
