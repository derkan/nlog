// +build windows

package console

import (
	"io"

	"golang.org/x/sys/windows"
)

// InitColor is needed to init windows console for ANSI colors
func InitColor(w io.Writer) {
	if f, ok := w.(file); ok {
		setConsoleMode(windows.Handle(f.Fd()), windows.ENABLE_VIRTUAL_TERMINAL_PROCESSING)
	}
}

// setConsoleMode sets the given flags on the given stream
func setConsoleMode(stdout windows.Handle, flags uint32) {
	var originalMode uint32
	windows.GetConsoleMode(stdout, &originalMode)
	windows.SetConsoleMode(stdout, originalMode|flags)
}
