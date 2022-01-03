package execTool

import (
	"log"
	"os"
	"path/filepath"
)

func ExecPath() string {
	file, err := os.Executable()
	if err != nil {
		log.Panicln(err)
	}
	return filepath.Dir(file)
}
