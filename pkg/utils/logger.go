package utils

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/fatih/color"
)

type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

type Logger struct {
	Level  Level
	Output io.Writer
}

var (
	DefaultLogger = &Logger{
		Level:  InfoLevel,
		Output: os.Stderr,
	}

	debugColor = color.New(color.FgHiBlack)
	infoColor  = color.New(color.FgCyan)
	warnColor  = color.New(color.FgYellow)
	errorColor = color.New(color.FgRed, color.Bold)
)

func SetLevel(level Level) {
	DefaultLogger.Level = level
}

func (l *Logger) log(level Level, c *color.Color, prefix, format string, args ...interface{}) {
	if level < l.Level {
		return
	}

	timestamp := time.Now().Format("15:04:05")
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(l.Output, "%s ", debugColor.Sprint(timestamp))
	c.Fprintf(l.Output, "[%s] ", prefix)
	fmt.Fprintf(l.Output, "%s\n", msg)
}

func Debug(format string, args ...interface{}) {
	DefaultLogger.log(DebugLevel, debugColor, "DEBUG", format, args...)
}

func Info(format string, args ...interface{}) {
	DefaultLogger.log(InfoLevel, infoColor, "INFO ", format, args...)
}

func Warn(format string, args ...interface{}) {
	DefaultLogger.log(WarnLevel, warnColor, "WARN ", format, args...)
}

func Error(format string, args ...interface{}) {
	DefaultLogger.log(ErrorLevel, errorColor, "ERROR", format, args...)
}
