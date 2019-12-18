package mylogger

import (
	"errors"
	"fmt"
	"path"
	"runtime"
	"strings"
	// "time"
)

type LogLevel uint16

type Logger interface {
	Debug(formt string, a ...interface{})
	Info(formt string, a ...interface{})
	Warning(formt string, a ...interface{})
	Error(formt string, a ...interface{})
	Trace(formt string, a ...interface{})
	Fatal(formt string, a ...interface{})
}

const (
	UNKNOWN LogLevel = iota
	DEBUG
	TRACE
	INFO
	WARNING
	ERROR
	FATAL
)

//parseLogLevel 根据字符串转换日志级别
func parseLogLevel(s string) (LogLevel, error) {
	s = strings.ToLower(s)
	switch s {
	case "debug":
		return DEBUG, nil
	case "trace":
		return TRACE, nil
	case "info":
		return INFO, nil
	case "error":
		return ERROR, nil
	case "warning":
		return WARNING, nil
	default:
		err := errors.New("无效的日志级别")
		return UNKNOWN, err
	}
}

//getLogString 将日志登记转为字符串
func getLogString(lv LogLevel) string {
	switch lv {
	case DEBUG:
		return "Debug"
	case TRACE:
		return "Trace"
	case INFO:
		return "Info"
	case WARNING:
		return "Warning"
	case ERROR:
		return "Error"
	case FATAL:
		return "Fatal"
	default:
		return "Debug"
	}
}

//getInfo...
func getInfo(skip int) (funcName, fileName string, lineNo int) {
	pc, file, lineNo, ok := runtime.Caller(skip)
	if !ok {
		fmt.Printf("runtime.caller() failed\n")
		return
	}
	funcName = runtime.FuncForPC(pc).Name()
	fileName = path.Base(file)
	funcName = strings.Split(funcName, ".")[1]
	return
}
