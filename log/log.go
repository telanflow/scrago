package log

import "fmt"

type Logs interface {
	Debug(msg interface{})
	Info(msg interface{})
	Notice(msg interface{})
	Error(msg interface{})
	Panic(msg interface{})
	DebugF(format string, v ...interface{})
	InfoF(format string, v ...interface{})
	NoticeF(format string, v ...interface{})
	ErrorF(format string, v ...interface{})
	PanicF(format string, v ...interface{})
}

var std = NewConsole()

func Debug(msg interface{}) {
	std.Debug(msg)
}

func Info(msg interface{}) {
	std.Info(msg)
}

func Notice(msg interface{}) {
	std.Notice(msg)
}

func Error(msg interface{}) {
	std.Error(msg)
}

func Panic(msg interface{}) {
	std.Panic(msg)
}

func DebugF(format string, v ...interface{}) {
	std.Debug(fmt.Sprintf(format, v...))
}

func InfoF(format string, v ...interface{}) {
	std.Info(fmt.Sprintf(format, v...))
}

func NoticeF(format string, v ...interface{}) {
	std.Notice(fmt.Sprintf(format, v...))
}

func ErrorF(format string, v ...interface{}) {
	std.Error(fmt.Sprintf(format, v...))
}

func PanicF(format string, v ...interface{}) {
	std.PanicF(fmt.Sprintf(format, v...))
}