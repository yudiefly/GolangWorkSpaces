package mylogger

import (
	"fmt"
	"time"
)

//往终端写入日志
// Logger 日志结构体
type ConsoleLogger struct {
	Level LogLevel
}

//NewConsoleLogger 初始化
func NewConsoleLogger(levelstr string) ConsoleLogger {
	level, err := parseLogLevel(levelstr)
	if err != nil {
		panic(err)
	}
	return ConsoleLogger{
		Level: level,
	}
}

//判断是否要记日志（等级）
func (c ConsoleLogger) enable(logLevel LogLevel) bool {
	return logLevel >= c.Level
}

//log 记日志的方法
func (c ConsoleLogger) log(lv LogLevel, format string, a ...interface{}) {
	if c.enable(lv) {
		msg := fmt.Sprintf(format, a...)
		now := time.Now()
		funcName, fileName, lineNo := getInfo(3)
		fmt.Printf("[%s] [%s] [%s:%s:%d]%s\n", now.Format("2006-01-02 15:04:05"), getLogString(lv), fileName, funcName, lineNo, msg)
	}
}

//Debug ...
func (c ConsoleLogger) Debug(format string, a ...interface{}) {
	c.log(DEBUG, format, a...)
}

//Info ...
func (c ConsoleLogger) Info(format string, a ...interface{}) {
	c.log(INFO, format, a...)
}

//Error ...
func (c ConsoleLogger) Error(format string, a ...interface{}) {

	c.log(ERROR, format, a...)

}

func (c ConsoleLogger) Warning(format string, a ...interface{}) {
	c.log(WARNING, format, a...)
}

//Fatal ...
func (c ConsoleLogger) Fatal(format string, a ...interface{}) {
	c.log(FATAL, format, a...)
}

//Trace ...
func (c ConsoleLogger) Trace(format string, a ...interface{}) {
	c.log(TRACE, format, a...)
}
