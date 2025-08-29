// 日志
package libsocket

type Logger interface {
	With(key string, value any) Logger
	Debug(args ...any)
	Debugf(format string, args ...any)
	Debugln(args ...any)
	Info(args ...any)
	Infof(format string, args ...any)
	Infoln(args ...any)
	Warn(args ...any)
	Warnf(format string, args ...any)
	Warnln(args ...any)
	Error(args ...any)
	Errorf(format string, args ...any)
	Errorln(args ...any)
}
