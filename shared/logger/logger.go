package logger

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/fatih/color"
)

type Logger struct {
	*log.Logger
	level LogLevel
}

type LogLevel struct {
	text color.Attribute
	tag  color.Attribute
}

var (
	INFO  = LogLevel{text: color.FgHiGreen, tag: color.BgHiGreen}
	DEBUG = LogLevel{text: color.FgCyan, tag: color.BgCyan}
	WARN  = LogLevel{text: color.FgYellow, tag: color.BgYellow}
	ERROR = LogLevel{text: color.FgRed, tag: color.BgRed}
)

func init() {
	color.NoColor = false
}

func New(level LogLevel) *Logger {
	return &Logger{
		Logger: log.New(os.Stdout, "", 0),
		level:  level,
	}
}

func getLocalTime() string {
	zone, _ := time.Now().Zone()      // Get the time zone name
	loc, _ := time.LoadLocation(zone) // Load the location using only the zone name
	now := time.Now().In(loc)

	return now.Format("01/02/2006 15:04:05") + fmt.Sprintf(".%06d", now.Nanosecond()/1000)
}

func (l *Logger) toOutput(format *string, v ...interface{}) {
	var message string

	timeString := getLocalTime()
	if format != nil {
		message = color.New(color.Attribute(l.level.text)).Sprintf(*format, v...)
	} else {
		message = color.New(color.Attribute(l.level.text)).Sprint(v...)
	}

	tag := color.New(color.Attribute(l.level.tag)).Sprintf(" %s ", l.levelString())
	output := fmt.Sprintf("API %s %s %s", timeString, tag, message)

	_ = l.Output(2, output)
}

func (l *Logger) Println(v ...interface{}) {
	l.toOutput(nil, v...)
}

func (l *Logger) Printf(format string, v ...interface{}) {
	l.toOutput(&format, v...)
}

func (l *Logger) Panicln(v ...interface{}) {
	l.toOutput(nil, v...)
	panic(v)
}

func (l *Logger) Panicf(format string, v ...interface{}) {
	l.toOutput(&format, v...)
	panic(v)
}

func (l *Logger) levelString() string {
	switch l.level.text {
	case INFO.text:
		return "INFO"
	case DEBUG.text:
		return "DEBUG"
	case WARN.text:
		return "WARN"
	case ERROR.text:
		return "ERROR"
	}
	return ""
}

var (
	Info  = New(LogLevel(INFO))
	Warn  = New(LogLevel(WARN))
	Debug = New(LogLevel(DEBUG))
	Error = New(LogLevel(ERROR))
)
