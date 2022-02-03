package pathTool

import (
	"os"
	"path/filepath"
)

// windows从运行文件目录获取工作目录路径
func SmartWd() string {
	file, _ := os.Executable()
	return filepath.Dir(file)
}

// windows转为绝对路径
func SmartAbs(path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(SmartWd(), path)
}
