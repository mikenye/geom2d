//go:build debug

package geom2d

import (
	"log"
	"os"
)

// Debug logger instance
var logger = log.New(os.Stderr, "[geom2d DEBUG] ", log.LstdFlags)

// Debug logs debug messages if the logger is enabled.
func logDebugf(format string, v ...interface{}) {
	logger.Printf(format, v...)
}
