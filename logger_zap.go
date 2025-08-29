// Package libsocket
// @description
// @author      梁志豪
// @datetime    2025/8/29 18:09
package libsocket

import "go.uber.org/zap"

// zapLogger 包装器
type zapLogger struct {
	sugar *zap.SugaredLogger
}

func NewLogger(sugar *zap.SugaredLogger) Logger {
	return &zapLogger{sugar: sugar}
}

// With 返回新的 Logger，绑定上下文字段
func (l *zapLogger) With(key string, value any) Logger {
	return &zapLogger{sugar: l.sugar.With(any(key), value)}
}

// Debug 系列
func (l *zapLogger) Debug(args ...any)                 { l.sugar.Debug(args...) }
func (l *zapLogger) Debugf(format string, args ...any) { l.sugar.Debugf(format, args...) }
func (l *zapLogger) Debugln(args ...any)               { l.sugar.Debug(args...) }

// Info 系列
func (l *zapLogger) Info(args ...any)                 { l.sugar.Info(args...) }
func (l *zapLogger) Infof(format string, args ...any) { l.sugar.Infof(format, args...) }
func (l *zapLogger) Infoln(args ...any)               { l.sugar.Info(args...) }

// Warn 系列
func (l *zapLogger) Warn(args ...any)                 { l.sugar.Warn(args...) }
func (l *zapLogger) Warnf(format string, args ...any) { l.sugar.Warnf(format, args...) }
func (l *zapLogger) Warnln(args ...any)               { l.sugar.Warn(args...) }

// Error 系列
func (l *zapLogger) Error(args ...any)                 { l.sugar.Error(args...) }
func (l *zapLogger) Errorf(format string, args ...any) { l.sugar.Errorf(format, args...) }
func (l *zapLogger) Errorln(args ...any)               { l.sugar.Error(args...) }
