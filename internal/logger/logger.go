package logger

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

var (
	infoLogger    *log.Logger
	errorLogger   *log.Logger
	warningLogger *log.Logger
	debugLogger   *log.Logger
)

func init() {
	// Create logs directory if not exists
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		log.Fatal("Failed to create logs directory:", err)
	}

	// Create or append to log file
	currentTime := time.Now()
	logFileName := filepath.Join(logsDir, fmt.Sprintf("%s.log", currentTime.Format("2006-01-02")))
	file, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("Failed to open log file:", err)
	}

	// Create multi-writer for both file and stdout
	multiWriter := io.MultiWriter(file, os.Stdout)

	// Initialize loggers with different prefixes
	infoLogger = log.New(multiWriter, "\033[32mINFO: \033[0m", log.Ldate|log.Ltime)
	errorLogger = log.New(multiWriter, "\033[31mERROR: \033[0m", log.Ldate|log.Ltime)
	warningLogger = log.New(multiWriter, "\033[33mWARNING: \033[0m", log.Ldate|log.Ltime)
	debugLogger = log.New(multiWriter, "\033[36mDEBUG: \033[0m", log.Ldate|log.Ltime)
}

func getFileAndLine() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "unknown"
		line = 0
	}
	return fmt.Sprintf("%s:%d", filepath.Base(file), line)
}

func Info(format string, v ...interface{}) {
	location := getFileAndLine()
	message := fmt.Sprintf(format, v...)
	infoLogger.Printf("[%s] %s", location, message)
}

func Error(format string, v ...interface{}) {
	location := getFileAndLine()
	message := fmt.Sprintf(format, v...)
	errorLogger.Printf("[%s] %s", location, message)
}

func Warning(format string, v ...interface{}) {
	location := getFileAndLine()
	message := fmt.Sprintf(format, v...)
	warningLogger.Printf("[%s] %s", location, message)
}

func Debug(format string, v ...interface{}) {
	location := getFileAndLine()
	message := fmt.Sprintf(format, v...)
	debugLogger.Printf("[%s] %s", location, message)
}
