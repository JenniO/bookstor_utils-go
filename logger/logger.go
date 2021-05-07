package logger

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"path"
	"runtime"
	"strings"
)

const (
	envLogLevel  = "LOG_LEVEL"
	envLogOutput = "LOG_OUTPUT"
)

var (
	log logger
)

type bookstoreLogger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
}

type logger struct {
	log *zap.Logger
}

func getLevel() zapcore.Level {
	switch strings.TrimSpace(strings.ToLower(os.Getenv(envLogLevel))) {
	case "debug":
		return zap.DebugLevel
	case "info":
		return zap.InfoLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}

func getOutput() string {
	output := strings.TrimSpace(os.Getenv(envLogLevel))
	if output == "" {
		return "stdout"
	}
	return output
}

func init() {
	logConfig := zap.Config{
		OutputPaths: []string{getOutput()},
		Level:       zap.NewAtomicLevelAt(getLevel()),
		Encoding:    "json",
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:   "msg",
			LevelKey:     "level",
			TimeKey:      "time",
			CallerKey:    "caller",
			EncodeTime:   zapcore.ISO8601TimeEncoder,
			EncodeLevel:  zapcore.LowercaseLevelEncoder,
			EncodeCaller: zapcore.ShortCallerEncoder,
		},
	}

	var err error
	if log.log, err = logConfig.Build(); err != nil {
		panic(err)
	}
}

func GetLogger() bookstoreLogger {
	return log
}

func (log logger) Print(v ...interface{}) {
	Info(fmt.Sprintf("%#v", v))
}

func (log logger) Printf(format string, v ...interface{}) {
	if len(v) == 0 {
		Info(format)
	} else {
		Info(fmt.Sprintf(format, v))
	}
}

func Info(msg string, tags ...zap.Field) {
	tags = append(tags, retrieveCallInfo())
	log.log.Info(msg, tags...)
	log.log.Sync()
}

func Error(msg string, err error, tags ...zap.Field) {
	tags = append(tags, zap.NamedError("error", err))
	tags = append(tags, retrieveCallInfo())
	log.log.Error(msg, tags...)
	log.log.Sync()
}

func retrieveCallInfo() zap.Field {
	programCounter, file, line, _ := runtime.Caller(1)
	functionCalled := strings.Split(runtime.FuncForPC(programCounter).Name(), ".")
	pl := len(functionCalled)
	funcName := functionCalled[pl-1]
	pkg := functionCalled[pl-2]
	pkg = strings.Trim(pkg, "()*")
	_, filename := path.Split(file)
	caller := fmt.Sprintf("%s/%s - %s : %d", pkg, filename, funcName, line)
	return zap.String("caller", caller)
}
