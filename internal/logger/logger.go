// Package logger реализует логгер (info, warn, error)
package logger

import (
	"log"
	"os"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
	warnLogger  *log.Logger
)

// Инициализация логгеров с соответствующими форматами вывода.
func init() {
	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	warnLogger = log.New(os.Stdout, "WARN: ", log.Ldate|log.Ltime|log.Lshortfile)
}

// Info записывает информационное сообщение
func Info(v ...interface{}) {
	infoLogger.Println(v...)
}

// Error записывает сообщение об ошибке
func Error(v ...interface{}) {
	errorLogger.Println(v...)
}

// Warn записывает предупреждающее сообщение
func Warn(v ...interface{}) {
	warnLogger.Println(v...)
}

// Info записывает форматируемое информационное сообщение
func Infof(format string, v ...interface{}) {
	infoLogger.Printf(format, v...)
}

// Error записывает форматируемое сообщение об ошибке
func Errorf(format string, v ...interface{}) {
	errorLogger.Printf(format, v...)
}

// Warn записывает форматируемое предупреждающее сообщение
func Warnf(format string, v ...interface{}) {
	warnLogger.Printf(format, v...)
}
