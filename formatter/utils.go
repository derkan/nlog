package formatter

import (
	"runtime"
	"strings"

	"github.com/derkan/nlog"
)

// GetFileLoc appends file location to log line
func GetFileLoc(pathStrip string, buff nlog.Buffer, callDepth int, quota bool) {
	var (
		filePath string
		lineNo   int
		ok       bool
	)

	if _, filePath, lineNo, ok = runtime.Caller(callDepth); !ok {
		filePath = "???"
		lineNo = 0
	}
	if pathStrip != "" {
		filePath = strings.TrimPrefix(filePath, pathStrip)
	}
	if quota {
		buff.AppendByte('"')
	}
	buff.AppendString(filePath, false).AppendByte(':').AppendInt(lineNo)
	if quota {
		buff.AppendByte('"')
	}
}
