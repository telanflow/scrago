package log

import (
	"log"
	"os"
	"fmt"
)

const (
	ColorBlack = iota + 30
	ColorRed
	ColorGreen
	ColorYellow
	ColorBlue
	ColorMagenta
	ColorCyan
	ColorWhite
)

type Console struct {
	log *log.Logger
}

func NewConsole() *Console {
	return &Console{
		log: log.New(os.Stdout, "", log.LstdFlags),
	}
}

func (self *Console) Debug(msg interface{}) {
	self.print(ColorGreen, msg)
}

func (self *Console) Info(msg interface{}) {
	self.print(ColorCyan, msg)
}

func (self *Console) Notice(msg interface{}) {
	self.print(ColorYellow, msg)
}

func (self *Console) Error(msg interface{}) {
	self.print(ColorRed, msg)
}

func (self *Console) Panic(msg interface{}) {
	self.print(ColorWhite, msg)
	panic(msg)
}

func (self *Console) DebugF(format string, v ...interface{}) {
	self.Debug(fmt.Sprintf(format, v...))
}

func (self *Console) InfoF(format string, v ...interface{}) {
	self.Info(fmt.Sprintf(format, v...))
}

func (self *Console) NoticeF(format string, v ...interface{}) {
	self.Notice(fmt.Sprintf(format, v...))
}

func (self *Console) ErrorF(format string, v ...interface{}) {
	self.Error(fmt.Sprintf(format, v...))
}

func (self *Console) PanicF(format string, v ...interface{}) {
	self.Panic(fmt.Sprintf(format, v...))
}

func (self *Console) print(color int, text interface{}) {
	self.log.Printf("\x1b[0;%dm%s\x1b[0m\n", color, text)
}