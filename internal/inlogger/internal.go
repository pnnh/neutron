package inlogger

import (
	"io"
	"path/filepath"
	"runtime"
	"strconv"

	"github.com/sirupsen/logrus"
)

var Logger = logrus.New()

// 初始化日志，默认什么也不会输出
func init() {
	Logger.SetOutput(io.Discard)
	Logger.SetLevel(logrus.PanicLevel)
	Logger.SetReportCaller(false)
	defaultFormatter := BuildFormatter()
	Logger.SetFormatter(defaultFormatter)
}

func BuildFormatter() *logrus.TextFormatter {
	return &logrus.TextFormatter{
		DisableColors:    false,
		DisableTimestamp: true,
		FullTimestamp:    true,
		TimestampFormat:  "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			return filepath.Base(f.Function), filepath.Base(f.File) + ":" + strconv.Itoa(f.Line)
		},
	}
}
