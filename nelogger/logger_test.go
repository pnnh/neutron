package nelogger

import (
	"github.com/pnnh/neutron/internal/inlogger"
)

func ExampleNEEnableLogger() {
	// 必须先启用日志输出，才能在控制台输出看到neutron包中的相关日志
	NEEnableLogger(true, InfoLevel, ShortFormat)

	// 仅内部能调用，此处为演示目的
	// 当调用其它neutron包中的函数时，这些函数会在内部调用inlogger相关函数来输出日志
	inlogger.Logger.Infoln("Hello")

	// Output:
	// level=info msg=Hello
}
