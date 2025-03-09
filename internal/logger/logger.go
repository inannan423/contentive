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
	auditLogger   *log.Logger
)

type UserType string

const (
	UserTypeAdmin UserType = "admin"
	UserTypeUser  UserType = "user"
)

func init() {
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, os.ModePerm); err != nil {
		log.Fatalf("failed to create logs directory: %v", err)
	}

	currentTime := time.Now()
	logFileName := filepath.Join(logsDir, fmt.Sprintf("log_%s.log", currentTime.Format("2006-01-02")))

	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open log file: %v", err)
	}

	infoLogger = log.New(logFile, "INFO: ", log.Ldate|log.Ltime)
	errorLogger = log.New(logFile, "ERROR: ", log.Ldate|log.Ltime)
	warningLogger = log.New(logFile, "WARNING: ", log.Ldate|log.Ltime)
	auditLogger = log.New(logFile, "AUDIT: ", log.Ldate|log.Ltime)

	// Output log messages to the console
	infoMultiWriter := io.MultiWriter(os.Stdout, logFile)
	warningMultiWriter := io.MultiWriter(os.Stdout, logFile)
	errorMultiWriter := io.MultiWriter(os.Stderr, logFile)
	auditMultiWriter := io.MultiWriter(os.Stdout, logFile)

	infoLogger = log.New(infoMultiWriter, "INFO: ", log.Ldate|log.Ltime)
	warningLogger = log.New(warningMultiWriter, "WARNING: ", log.Ldate|log.Ltime)
	errorLogger = log.New(errorMultiWriter, "ERROR: ", log.Ldate|log.Ltime)
	auditLogger = log.New(auditMultiWriter, "AUDIT: ", log.Ldate|log.Ltime)
}

// getFileAndLine returns the file name and line number of the caller
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

func Warning(format string, v ...interface{}) {
	location := getFileAndLine()
	message := fmt.Sprintf(format, v...)
	warningLogger.Printf("[%s] %s", location, message)
}

func Error(format string, v ...interface{}) {
	location := getFileAndLine()
	message := fmt.Sprintf(format, v...)
	errorLogger.Printf("[%s] %s", location, message)
}

// AdminAction logs admin actions
func AdminAction(userID uint, username string, action string, details string) {
	location := getFileAndLine()
	auditLogger.Printf("[%s] [AdminUser:%d:%s] %s - %s", location, userID, username, action, details)
}

// APIAction logs API actions
func APIAction(userID string, username string, action string, details string) {
	location := getFileAndLine()
	auditLogger.Printf("[%s] [APIUser:%s:%s] %s - %s", location, userID, username, action, details)
}

// GeneralAction logs general actions
func GeneralAction(details string) {
	location := getFileAndLine()
	auditLogger.Printf("[%s][System] %s", location, details)
}
