package nelogger

import (
	"io"
	"os"

	"github.com/pnnh/neutron/internal/inlogger"

	"github.com/sirupsen/logrus"
)

type NELogLevel uint32

const (
	PanicLevel NELogLevel = iota
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

type NELogFormat uint32

const (
	ShortFormat NELogFormat = iota
	FullFormat
)

// NEEnableLogger 启用neutron库中的日志输出功能，需要同时指定一个日志级别，和日志格式
func NEEnableLogger(enable bool, level NELogLevel, format NELogFormat) {
	if enable {
		inlogger.Logger.SetOutput(os.Stdout)
		NESetLevel(level)
		SetFormat(format)
	} else {
		inlogger.Logger.SetOutput(io.Discard)
	}
}

// SetFormat 设置日志格式，目前支持长短两种格式
func SetFormat(format NELogFormat) {
	formatter := inlogger.BuildFormatter()
	if format == ShortFormat {
		formatter.DisableTimestamp = true
	} else {
		formatter.DisableTimestamp = false
	}
	inlogger.Logger.SetFormatter(formatter)
}

// NESetLevel 设置neutron库内部日志的输出级别
func NESetLevel(level NELogLevel) {
	switch level {
	case PanicLevel:
		inlogger.Logger.SetLevel(logrus.PanicLevel)
	case FatalLevel:
		inlogger.Logger.SetLevel(logrus.FatalLevel)
	case ErrorLevel:
		inlogger.Logger.SetLevel(logrus.ErrorLevel)
	case WarnLevel:
		inlogger.Logger.SetLevel(logrus.WarnLevel)
	case InfoLevel:
		inlogger.Logger.SetLevel(logrus.InfoLevel)
	case DebugLevel:
		inlogger.Logger.SetLevel(logrus.DebugLevel)
	case TraceLevel:
		inlogger.Logger.SetLevel(logrus.TraceLevel)
	default:
		inlogger.Logger.SetLevel(logrus.InfoLevel)
	}
}

func NELogSetReportCaller(enabled bool) {
	inlogger.Logger.SetReportCaller(enabled)
}
