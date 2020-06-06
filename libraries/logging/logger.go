package logging

import (
	"fmt"
	"sync"
	"time"
)

// Logging Levels
const (
	LogLevelError = iota
	LogLevelInformation
	LogLevelWarning
	LogLevelDebug
)

func init() {
	// Intantiates the default logger with a stdout output.
	defaultLogger = &Logger{
		mu:      &sync.RWMutex{},
		outputs: []Output{StdOut()},
	}
}

var defaultLogger *Logger

// Output is a heigh-level interface used to log messages by Logger.
type Output interface {
	Write(text string)
}

// Logger is struct which contains methods used to log messages to a number
// of outputs, with supports for logging levels.
type Logger struct {
	lvl     int
	mu      *sync.RWMutex
	outputs []Output
}

// Information logs a message to the default logger.
func Information(message string) {
	defaultLogger.Information(message)
}

// Informationf logs a message to the default logger using fmt.Sprintf.
func Informationf(format string, v ...interface{}) {
	defaultLogger.Informationf(format, v...)
}

// Warning logs a message to the default logger.
func Warning(message string) {
	defaultLogger.Warning(message)
}

// Warningf logs a message to the default logger using fmt.Sprintf.
func Warningf(format string, v ...interface{}) {
	defaultLogger.Warningf(format, v...)
}

// Debug logs a message to the default logger.
func Debug(message string) {
	defaultLogger.Debug(message)
}

// Debugf logs a message to the default logger using fmt.Sprintf.
func Debugf(format string, v ...interface{}) {
	defaultLogger.Debugf(format, v...)
}

// Error is used to log an error to the default logger.
func Error(err error) {
	defaultLogger.Error(err)
}

// Errorf logs a message to the default logger using fmt.Sprintf.
func Errorf(format string, v ...interface{}) {
	defaultLogger.Errorf(format, v...)
}

func WithOutput(o Output) *Logger {
	return defaultLogger.WithOutput(o)
}

func SetLogLevel(lvl int) *Logger {
	return defaultLogger.SetLogLevel(lvl)
}

// ResetOutputs clears the output options of the default logger.
func ResetOutputs() *Logger {
	return defaultLogger.ResetOutputs()
}

func (l *Logger) Information(message string) {
	l.output(LogLevelInformation, formater(LogLevelInformation, message))
}

func (l *Logger) Informationf(format string, v ...interface{}) {
	l.output(LogLevelInformation, formater(LogLevelInformation, format, v...))
}

func (l *Logger) Warning(message string) {
	l.output(LogLevelWarning, formater(LogLevelWarning, message))
}

func (l *Logger) Warningf(format string, v ...interface{}) {
	l.output(LogLevelWarning, formater(LogLevelWarning, format, v...))
}

func (l *Logger) Debug(message string) {
	l.output(LogLevelDebug, formater(LogLevelDebug, message))
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.output(LogLevelDebug, formater(LogLevelDebug, format, v...))
}

func (l *Logger) Error(err error) {
	l.Errorf("%v", err)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.output(LogLevelError, formater(LogLevelError, format, v...))
}

// WithOutput adds an output option to the Logger.
func (l *Logger) WithOutput(o Output) *Logger {
	l.mu.RLock()
	defer l.mu.RUnlock()

	l.outputs = append(l.outputs, o)

	return l
}

// SetLogLevel sets the logging level of the Logger.
func (l *Logger) SetLogLevel(lvl int) *Logger {
	l.lvl = lvl

	return l
}

// ResetOutputs clear all output options for the Logger.
func (l *Logger) ResetOutputs() *Logger {
	l.mu.RLock()
	defer l.mu.RUnlock()

	l.outputs = []Output{}

	return l
}

// writes a message to the Logger's outputs, if the log level is enabled.
func (l *Logger) output(lvl int, message string) {
	if !l.isEnabled(lvl) {
		return
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	for _, o := range l.outputs {
		o.Write(message)
	}
}

// formats a string for the given lvl.
func formater(lvl int, format string, v ...interface{}) string {
	var prefix string
	switch lvl {
	case LogLevelDebug:
		prefix = "[DEBUG]"
		break
	case LogLevelError:
		prefix = "[ERROR]"
		break
	case LogLevelWarning:
		prefix = "[WARN]"
		break
	case LogLevelInformation:
		prefix = "[INFO]"
		break
	}

	t := time.Now().UTC()
	tf := fmt.Sprintf("[%d/%d/%d %d:%d:%d]",
		t.Year(),
		t.Month(),
		t.Day(),
		t.Hour(),
		t.Minute(),
		t.Second())
	msg := fmt.Sprintf(format, v...)

	return fmt.Sprintf("%s%s: %s", tf, prefix, msg)
}

func (l *Logger) isEnabled(lvl int) bool {
	if l.lvl == 0 {
		return true
	}

	return lvl <= l.lvl
}
