package pathTool

import (
	"os"
	"path/filepath"
)

// linux正常获取工作目录路径
func SmartWd() string {
	wd, _ := os.Getwd()
	return wd
}

// linux转为绝对路径
func SmartAbs(path string) string {
	abs, _ := filepath.Abs(path)
	return abs
}
