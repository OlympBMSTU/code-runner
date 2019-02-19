package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/OlympBMSTU/code-runner/config"
)

/// Trash
var loggerInstance *Logger = nil

func InitLogger(cfg config.Config) error {
	// if logger == nil {
	logger, err := New(cfg)
	loggerInstance = logger
	return err
	// }
	// return
}

// cant as ILogger ...
func GetLogger() *Logger { //cfg config.Config) (*Logger, error) {
	// if logger == nil {
	// logger, err := New(cfg)
	// return logger, err
	// }

	return loggerInstance
}

////

// LOGGER constants for create logger
const (
	CONSOLE      = "console"
	FILE         = "file"
	BOTH         = "console_and_file"
	FULL_LOG     = 2
	EXTENDED_LOG = 1
	MIN_LOG      = 0
)

const (
	SKIP_FUNCTION_STACK = 2
	LINE_TERMINATOR     = "\n"
)

// logger logs sended data
// if level 0 - only requred parameters
// if level 1 - all parameters
// if level 2 appends string that tell where log writed line func file
type ILogger interface {
	Info(string, ...interface{})
	Warn(string, ...interface{})
	Error(string, error, ...interface{})
}

type Logger struct {
	level  int
	writer LogWriter
}

func New(config config.Config) (*Logger, error) {
	writer := LogWriter{[]*os.File{os.Stdout}}
	if config.GetLoggerType() == FILE {
		file, err := os.OpenFile(config.GetLoggerPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return nil, err
		}
		writer = LogWriter{[]*os.File{file}}
	}

	if config.GetLoggerType() == BOTH {
		file, err := os.OpenFile(config.GetLoggerPath(), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return nil, err
		}
		writer = LogWriter{[]*os.File{os.Stdout, file}}
	}

	return &Logger{
		level:  config.GetLoggerLevel(),
		writer: writer,
	}, nil
}

func (logger Logger) getLogString(logData string, args ...interface{}) string {
	if logger.level == 0 {
		return logData
	} else {
		for _, arg := range args {
			if reflect.TypeOf(arg).Kind() != reflect.String {
				bytes, _ := json.Marshal(arg)
				logData += string(bytes)
				logData += " "
			} else {
				logData += arg.(string)
				logData += " "
			}
		}
		return logData
	}
}

func (logger Logger) Info(info string, args ...interface{}) {
	if logger.level == 2 {
		additionalInfo := getFuncInfo()
		args = append(args, additionalInfo)
	}
	current := time.Now()
	logStr := "Time: " + fmt.Sprintf("%s\n", current.String())
	logStr += "INFO: " + info + LINE_TERMINATOR
	str := logger.getLogString(logStr, args...)
	logger.writer.Write(str)
}

func (logger Logger) Warn(info string, args ...interface{}) {
	if logger.level == 2 {
		additionalInfo := getFuncInfo()
		args = append(args, additionalInfo)
	}
	current := time.Now()
	logStr := "Time: " + fmt.Sprintf("%s\n", current.String())
	logStr += "WARN: " + info + LINE_TERMINATOR
	str := logger.getLogString(logStr, args...)
	logger.writer.Write(str)
}

func (logger Logger) Error(msg string, err error, args ...interface{}) {
	if logger.level == 2 {
		additionalInfo := getFuncInfo()
		args = append(args, additionalInfo)
	}
	errMsg := ""
	if err != nil {
		errMsg = err.Error() + "\n"
	}
	bytes, _ := json.Marshal(err)
	errMsg += string(bytes)
	current := time.Now()
	logStr := "Time: " + fmt.Sprintf("%s\n", current.String())
	logStr += "ERR: " + msg + LINE_TERMINATOR + errMsg + LINE_TERMINATOR
	str := logger.getLogString(logStr, args...)
	logger.writer.Write(str)
}

func getFuncInfo() string {
	additionalInfo := ""
	pc, file, line, ok := runtime.Caller(SKIP_FUNCTION_STACK) // 2
	if !ok {
		file = "?"
		line = 0
	}

	fn := runtime.FuncForPC(pc)
	var fnName string
	if fn == nil {
		fnName = "?()"
	} else {
		dotName := filepath.Ext(fn.Name())
		fnName = strings.TrimLeft(dotName, ".") + "()"
	}

	additionalInfo = fmt.Sprintf("File: %s Func: %s Line: %d", filepath.Base(file), fnName, line)
	return additionalInfo + LINE_TERMINATOR
}
