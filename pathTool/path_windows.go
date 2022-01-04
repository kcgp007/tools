package pathTool

import (
	"os"
	"path/filepath"
)

func SmartWd() string {
	file, _ := os.Executable()
	return filepath.Dir(file)
}

func SmartAbs(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(SmartWd(), path)
}
