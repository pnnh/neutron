package nelogger

import (
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
