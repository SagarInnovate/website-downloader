package utils

import (
	"log"
	"os"
)

var (
	InfoLogger  *log.Logger
	ErrorLogger *log.Logger
)

func init() {
	InfoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	ErrorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// LogInfo logs informational messages
func LogInfo(format string, v ...interface{}) {
	InfoLogger.Printf(format, v...)
}

// LogError logs error messages
func LogError(format string, v ...interface{}) {
	ErrorLogger.Printf(format, v...)
}
