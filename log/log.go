package log

import (
	"os"
	"strings"

	"github.com/op/go-logging"
)

var (
	logger *logging.Logger
	module string
)

var backend logging.Backend
var leveledBackend logging.LeveledBackend

const DefaultFormat = "%{time:2006-01-02 15:04:05.000} | %{message}"

func init() {
	module = os.Args[0]
	Configure(DefaultFormat, logging.INFO)
}

func Configure(format string, level logging.Level) {
	logger = NewLogger(format, level)
}

func NewLogger(format string, level logging.Level) (logger *logging.Logger) {
	logger = logging.MustGetLogger(module)
	formatter := logging.MustStringFormatter(format)

	backend = logging.NewLogBackend(os.Stdout, "", 0)
	backFormatter := logging.NewBackendFormatter(backend, formatter)
	leveledBackend = logging.AddModuleLevel(backFormatter)
	leveledBackend.SetLevel(level, module)

	logger.SetBackend(leveledBackend)
	return
}

func SetLevel(level logging.Level) {
	leveledBackend.SetLevel(level, module)
}

func SetLevelFromStr(str string) {
	levelStr := strings.ToUpper(str)
	var level logging.Level
	switch levelStr {
	case "CRITICAL":
		level = logging.CRITICAL
	case "ERROR":
		level = logging.ERROR
	case "WARNING":
		level = logging.WARNING
	case "NOTICE":
		level = logging.NOTICE
	case "INFO":
		level = logging.INFO
	case "DEBUG":
		level = logging.DEBUG
	}

	SetLevel(level)
}

func SetFormat(format string) {
	level := logging.GetLevel(module)
	logging.Reset()
	NewLogger(format, level)
}

func Debugf(format string, args ...interface{}) {
	logger.Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	logger.Infof(format, args...)
}

func Noticef(format string, args ...interface{}) {
	logger.Noticef(format, args...)
}

func Warnf(format string, args ...interface{}) {
	logger.Warningf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logger.Errorf(format, args...)
}

func Criticalf(format string, args ...interface{}) {
	logger.Criticalf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	logger.Panicf(format, args...)
}
