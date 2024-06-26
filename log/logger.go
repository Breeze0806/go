package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"sync"
)

// Logger 用于打印调试日志
type Logger interface {
	Errorf(format string, v ...interface{}) //错误日志打印
	Warnf(format string, v ...interface{})  //告警日志打印
	Infof(format string, v ...interface{})  //进程日志打印
	Debugf(format string, v ...interface{}) //调试日志打印
	Print(args ...interface{})              //打印错误日志
	Printf(format string, v ...interface{}) //打印错误日志
}

// Level 日志级别, 为调试/信息/错误
type Level uint8

// 日志级别
const (
	DebugLevel Level = iota //调试
	InfoLevel               //信息
	WarnLevel               //告警
	ErrorLevel              //错误
)

type defaultLogger struct {
	level  Level
	logger *log.Logger
}

func newNilLogger() Logger {
	d := &defaultLogger{
		level:  ErrorLevel,
		logger: log.New(os.Stderr, "[log]", log.Lmicroseconds|log.LstdFlags|log.Llongfile),
	}
	return d
}

// NewDefaultLogger 生成一个日志打印Logger，level可以是DebugLevel，InfoLevel，ErrorLevel
func NewDefaultLogger(writer io.Writer, level Level, prefix string) Logger {
	d := &defaultLogger{
		level:  level,
		logger: log.New(writer, prefix, log.Lmicroseconds|log.LstdFlags|log.Llongfile),
	}
	return d
}

// Errorf 错误日志打印
func (d *defaultLogger) Errorf(format string, args ...interface{}) {
	if d.level <= ErrorLevel {
		b := &strings.Builder{}
		b.WriteString("[ERROR] ")
		b.WriteString(fmt.Sprintf(format, args...))
		d.logger.Output(2, b.String())
	}
}

// Warnf 错误日志打印
func (d *defaultLogger) Warnf(format string, args ...interface{}) {
	if d.level <= WarnLevel {
		b := &strings.Builder{}
		b.WriteString("[WARN] ")
		b.WriteString(fmt.Sprintf(format, args...))
		d.logger.Output(2, b.String())
	}
}

// Infof 进程日志打印
func (d *defaultLogger) Infof(format string, args ...interface{}) {
	if d.level <= InfoLevel {
		b := &strings.Builder{}
		b.WriteString("[INFO] ")
		b.WriteString(fmt.Sprintf(format, args...))
		d.logger.Output(2, b.String())
	}
}

// Debugf 进程日志打印
func (d *defaultLogger) Debugf(format string, args ...interface{}) {
	if d.level <= DebugLevel {
		b := &strings.Builder{}
		b.WriteString("[DEBUG] ")
		b.WriteString(fmt.Sprintf(format, args...))
		d.logger.Output(2, b.String())
	}
}

// Print 日志打印
func (d *defaultLogger) Print(args ...interface{}) {
	d.logger.Output(2, fmt.Sprint(args...))
}

// Printf 日志打印
func (d *defaultLogger) Printf(format string, v ...interface{}) {
	d.logger.Output(2, fmt.Sprintf(format, v...))
}

var (
	lw    = loggerWrapper{l: newNilLogger()}
	lfuns = loggerFuncs{funcs: make([]func(), 0, 10)}
)

type loggerFuncs struct {
	sync.Mutex
	funcs []func()
}

func (l *loggerFuncs) append(f func()) {
	l.Lock()
	defer l.Unlock()
	l.funcs = append(l.funcs, f)
}

func (l *loggerFuncs) doAllFuns() {
	l.Lock()
	defer l.Unlock()
	for _, f := range l.funcs {
		if f != nil {
			f()
		}
	}
}

type loggerWrapper struct {
	l  Logger
	mu sync.RWMutex
}

func (l *loggerWrapper) setLogger(logger Logger) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.l = logger
}

func (l *loggerWrapper) logger() Logger {
	l.mu.RLock()
	defer l.mu.RUnlock()
	return l.l
}

// SetLogger 设置一个符合Logger日志来打印调试信息，并且执行所有的日志初始化函数
func SetLogger(logger Logger) {
	lw.setLogger(logger)
	lfuns.doAllFuns()
}

// GetLogger 获取日志答应句柄
func GetLogger() Logger {
	return lw.logger()
}

// RegisterInitFuncs 注册获取初始化函数
func RegisterInitFuncs(f func()) {
	lfuns.append(f)
}
