package dglogger

import dgctx "github.com/darwinOrg/go-common/context"

var GlobalDgLogger = DefaultDgLogger()

func Debugf(ctx *dgctx.DgContext, format string, args ...any) {
	GlobalDgLogger.Debugf(ctx, format, args...)
}

func Infof(ctx *dgctx.DgContext, format string, args ...any) {
	GlobalDgLogger.Infof(ctx, format, args...)
}

func Warnf(ctx *dgctx.DgContext, format string, args ...any) {
	GlobalDgLogger.Warnf(ctx, format, args...)
}

func Errorf(ctx *dgctx.DgContext, format string, args ...any) {
	GlobalDgLogger.Errorf(ctx, format, args...)
}

func Fatalf(ctx *dgctx.DgContext, format string, args ...any) {
	GlobalDgLogger.Fatalf(ctx, format, args...)
}

func Panicf(ctx *dgctx.DgContext, format string, args ...any) {
	GlobalDgLogger.Panicf(ctx, format, args...)
}

func Debug(ctx *dgctx.DgContext, args ...any) {
	GlobalDgLogger.Debug(ctx, args...)
}

func Info(ctx *dgctx.DgContext, args ...any) {
	GlobalDgLogger.Info(ctx, args...)
}

func Warn(ctx *dgctx.DgContext, args ...any) {
	GlobalDgLogger.Warn(ctx, args...)
}

func Error(ctx *dgctx.DgContext, args ...any) {
	GlobalDgLogger.Error(ctx, args...)
}

func Fatal(ctx *dgctx.DgContext, args ...any) {
	GlobalDgLogger.Fatal(ctx, args...)
}

func Panic(ctx *dgctx.DgContext, args ...any) {
	GlobalDgLogger.Panic(ctx, args...)
}

func Debugln(ctx *dgctx.DgContext, args ...any) {
	GlobalDgLogger.Debugln(ctx, args...)
}

func Infoln(ctx *dgctx.DgContext, args ...any) {
	GlobalDgLogger.Infoln(ctx, args...)
}

func Warnln(ctx *dgctx.DgContext, args ...any) {
	GlobalDgLogger.Warnln(ctx, args...)
}

func Errorln(ctx *dgctx.DgContext, args ...any) {
	GlobalDgLogger.Errorln(ctx, args...)
}

func Fatalln(ctx *dgctx.DgContext, args ...any) {
	GlobalDgLogger.Fatalln(ctx, args...)
}

func Panicln(ctx *dgctx.DgContext, args ...any) {
	GlobalDgLogger.Panicln(ctx, args...)
}
