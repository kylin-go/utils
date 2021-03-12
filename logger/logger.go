package logger

import (
	"io"
	"log"
	"os"
	"runtime"
)

const (
	erro  = 3
	warn  = 2
	info  = 1
	debug = 0
)

type logger struct {
	info   *log.Logger
	warn   *log.Logger
	error  *log.Logger
	debug  *log.Logger
	except *log.Logger
	writer *io.Writer
	level  int
}

func (l *logger) init() {
	l.info = log.New(*l.writer, "[INFO ]: ", log.Ldate|log.Ltime)
	l.warn = log.New(*l.writer, "[WARN ]: ", log.Ldate|log.Ltime)
	l.error = log.New(*l.writer, "[ERROR]: ", log.Ldate|log.Ltime)
	l.debug = log.New(*l.writer, "[DEBUG]: ", log.Ldate|log.Ltime)
	l.except = log.New(*l.writer, "[EXECP]: ", log.Ldate|log.Ltime)
}

func (l *logger) Info(v ...interface{}) {
	if l.level <= info {
		l.info.Println(v...)
	}
}

func (l *logger) Infof(format string, v ...interface{}) {
	if l.level <= info {
		l.info.Printf(format, v...)
	}
}

func (l *logger) Warn(v ...interface{}) {
	if l.level <= warn {
		l.warn.Println(v...)
	}
}

func (l *logger) Warnf(format string, v ...interface{}) {
	if l.level <= warn {
		l.warn.Printf(format, v...)
	}
}

func (l *logger) Error(v ...interface{}) {
	if l.level <= erro {
		l.error.Println(v...)
	}
}

func (l *logger) Errorf(format string, v ...interface{}) {
	if l.level <= erro {
		l.error.Printf(format, v...)
	}
}

func (l *logger) Debug(v ...interface{}) {
	if l.level <= debug {
		l.debug.Println(v...)
	}
}

func (l *logger) Debugf(format string, v ...interface{}) {
	if l.level <= debug {
		l.debug.Printf(format, v...)
	}
}

func (l *logger) Exception(v ...interface{}) {
	l.except.Println(v...)
	var buf [4096]byte
	n := runtime.Stack(buf[:], true)
	l.except.Println(string(buf[:n]))
}

func (l *logger) Exceptionf(format string, v ...interface{}) {
	l.except.Printf(format, v...)
	var buf [4096]byte
	n := runtime.Stack(buf[:], true)
	l.except.Println(string(buf[:n]))
}

func NewLogger(logLevel, logPath string, console bool) *logger {
	var logWriter []io.Writer
	var levelMap = map[string]int{
		"debug": 0,
		"info":  1,
		"warn":  2,
		"error": 3,
	}
	var level, ok = levelMap[logLevel]
	if !ok {
		w := io.MultiWriter(os.Stderr)
		l := log.New(w, "[WARN ]: ", log.Ldate|log.Ltime)
		l.Println("level 值错误, 使用默认值info, 正确值为\"debug, info, warn, error\"")
		level = 1
	}
	if console {
		logWriter = append(logWriter, os.Stdout)
	}
	f, err := os.OpenFile(logPath, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err == nil {
		logWriter = append(logWriter, f)
	}
	writer := io.MultiWriter(logWriter...)
	l := logger{writer: &writer, level: level}
	l.init()
	return &l
}
