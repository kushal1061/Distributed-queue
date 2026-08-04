package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger

func InitLogger() {

	allLogFilePath := "logs/app.log" //this will store all logs for debugging
	err := os.MkdirAll("logs", os.ModePerm)
	if err != nil {
		panic("Failed to create log directory " + err.Error())
	}
	allLogFile, err := os.OpenFile(allLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("Failed to open log file: " + err.Error())
	}
	errorLogFile, err := os.OpenFile("logs/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic("Failed to open error.log file: " + err.Error())
	}
	// Encoders
	allLogsEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	errorLogsEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	// Log cores
	allLogsCore := zapcore.NewCore(allLogsEncoder, zapcore.AddSync(allLogFile), zapcore.DebugLevel)       // All logs
	errorLogsCore := zapcore.NewCore(errorLogsEncoder, zapcore.AddSync(errorLogFile), zapcore.ErrorLevel) // Only errors
	consoleCore := zapcore.NewCore(consoleEncoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel)        // Console output

	// Combine cores (console + all logs + error logs)
	core := zapcore.NewTee(allLogsCore, errorLogsCore, consoleCore)

	// Build logger
	Logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	defer allLogFile.Close()
	defer errorLogFile.Close()

	//explore lumberjack to rotate logs, prevent app.log from getting too large....
}
func GetLogger() *zap.Logger {
	if Logger == nil {
		InitLogger()
	}
	return Logger
}
