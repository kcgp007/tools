package tools

import (
	"syscall"

	"github.com/fatih/color"
	sequences "github.com/konsorten/go-windows-terminal-sequences"
)

func init() {
	if err := sequences.EnableVirtualTerminalProcessing(syscall.Stdout, true); err != nil {
		color.NoColor = false
	}
}
