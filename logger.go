package dglogger

import (
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"

	"github.com/darwinOrg/go-common/constants"
	dgctx "github.com/darwinOrg/go-common/context"
	dgsys "github.com/darwinOrg/go-common/sys"
	"github.com/darwinOrg/go-common/utils"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// PanicLevel level, highest level of severity. Logs and then calls panic with the
	// message passed to Debug, Info, ...
	PanicLevel = "panic"
	// FatalLevel level. Logs and then calls `logger.Exit(1)`. It will exit even if the
	// logging level is set to Panic.
	FatalLevel = "fatal"
	// ErrorLevel level. Logs. Used for errors that should definitely be noted.
	// Commonly used for hooks to send errors to an error tracking service.
	ErrorLevel = "error"
	// WarnLevel level. Non-critical entries that deserve eyes.
	WarnLevel = "warn"
	// InfoLevel level. General operational entries about what's going on inside the
	// application.
	InfoLevel = "info"
	// DebugLevel level. Usually only enabled when debugging. Very verbose logging.
	DebugLevel = "debug"
)

const (
	DefaultTimestampFormat = "2006-01-02 15:04:05.999999"
	DefaultFilename        = "app.log" // 日志文件路径
	DefaultMaxSize         = 100       // 每个日志文件的最大尺寸（MB）
	DefaultMaxBackups      = 10        // 保留旧日志文件的最大数量
	DefaultMaxAge          = 30        // 保留旧日志文件的最大天数
	DefaultCompress        = true      // 是否压缩/归档旧的日志文件
	extraFieldsKey         = "extraLogFields"
)

type DgLogger struct {
	log *zap.Logger
}

func DefaultDgLogger() *DgLogger {
	return NewDgLogger(getDefaultLogLevel(), DefaultTimestampFormat, os.Stdout)
}

func DefaultRotatedLogger() *DgLogger {
	return NewDgLogger(getDefaultLogLevel(), DefaultTimestampFormat, buildDefaultRotatedLogWriter())
}

func DefaultMultiWriterLogger() *DgLogger {
	return NewDgLogger(getDefaultLogLevel(), DefaultTimestampFormat, io.MultiWriter(os.Stdout, buildDefaultRotatedLogWriter()))
}

func NewDgLogger(level string, timestampFormat string, out io.Writer) *DgLogger {
	// 创建 zap 配置
	config := zap.NewProductionEncoderConfig()
	config.EncodeTime = zapcore.TimeEncoderOfLayout(timestampFormat)
	config.EncodeLevel = zapcore.CapitalLevelEncoder

	// 创建 encoder
	encoder := zapcore.NewConsoleEncoder(config)

	// 创建 core
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(out),
		parseLevel(level),
	)

	// 创建 logger
	logger := zap.New(core, zap.AddCallerSkip(1))

	return &DgLogger{log: logger}
}

func getDefaultLogLevel() string {
	return utils.IfReturn(dgsys.IsProd(), InfoLevel, DebugLevel)
}

func buildDefaultRotatedLogWriter() io.Writer {
	return &lumberjack.Logger{
		Filename:   DefaultFilename,
		MaxSize:    DefaultMaxSize,
		MaxBackups: DefaultMaxBackups,
		MaxAge:     DefaultMaxAge,
		Compress:   DefaultCompress,
	}
}

func (dl *DgLogger) Debugf(ctx *dgctx.DgContext, format string, args ...any) {
	if dl.log.Core().Enabled(zap.DebugLevel) {
		dl.withFields(ctx, nil, false).Debug(fmt.Sprintf(format, args...))
	}
}

func (dl *DgLogger) Infof(ctx *dgctx.DgContext, format string, args ...any) {
	if dl.log.Core().Enabled(zap.InfoLevel) {
		dl.withFields(ctx, nil, false).Info(fmt.Sprintf(format, args...))
	}
}

func (dl *DgLogger) Warnf(ctx *dgctx.DgContext, format string, args ...any) {
	if dl.log.Core().Enabled(zap.WarnLevel) {
		dl.withFields(ctx, nil, false).Warn(fmt.Sprintf(format, args...))
	}
}

func (dl *DgLogger) Errorf(ctx *dgctx.DgContext, format string, args ...any) {
	if dl.log.Core().Enabled(zap.ErrorLevel) {
		dl.withFields(ctx, nil, true).Error(fmt.Sprintf(format, args...))
	}
}

func (dl *DgLogger) Fatalf(ctx *dgctx.DgContext, format string, args ...any) {
	dl.withFields(ctx, nil, true).Fatal(fmt.Sprintf(format, args...))
}

func (dl *DgLogger) Panicf(ctx *dgctx.DgContext, format string, args ...any) {
	dl.withFields(ctx, nil, true).Panic(fmt.Sprintf(format, args...))
}

func (dl *DgLogger) Debug(ctx *dgctx.DgContext, args ...any) {
	if dl.log.Core().Enabled(zap.DebugLevel) {
		dl.withFields(ctx, nil, false).Debug(fmt.Sprint(args...))
	}
}

func (dl *DgLogger) Info(ctx *dgctx.DgContext, args ...any) {
	if dl.log.Core().Enabled(zap.InfoLevel) {
		dl.withFields(ctx, nil, false).Info(fmt.Sprint(args...))
	}
}

func (dl *DgLogger) Warn(ctx *dgctx.DgContext, args ...any) {
	if dl.log.Core().Enabled(zap.WarnLevel) {
		dl.withFields(ctx, nil, false).Warn(fmt.Sprint(args...))
	}
}

func (dl *DgLogger) Error(ctx *dgctx.DgContext, args ...any) {
	if dl.log.Core().Enabled(zap.ErrorLevel) {
		dl.withFields(ctx, nil, true).Error(fmt.Sprint(args...))
	}
}

func (dl *DgLogger) Fatal(ctx *dgctx.DgContext, args ...any) {
	dl.withFields(ctx, nil, true).Fatal(fmt.Sprint(args...))
}

func (dl *DgLogger) Panic(ctx *dgctx.DgContext, args ...any) {
	dl.withFields(ctx, nil, true).Panic(fmt.Sprint(args...))
}

func (dl *DgLogger) Debugln(ctx *dgctx.DgContext, args ...any) {
	if dl.log.Core().Enabled(zap.DebugLevel) {
		dl.withFields(ctx, nil, false).Debug(fmt.Sprintln(args...))
	}
}

func (dl *DgLogger) Infoln(ctx *dgctx.DgContext, args ...any) {
	if dl.log.Core().Enabled(zap.InfoLevel) {
		dl.withFields(ctx, nil, false).Info(fmt.Sprintln(args...))
	}
}

func (dl *DgLogger) Warnln(ctx *dgctx.DgContext, args ...any) {
	if dl.log.Core().Enabled(zap.WarnLevel) {
		dl.withFields(ctx, nil, false).Warn(fmt.Sprintln(args...))
	}
}

func (dl *DgLogger) Errorln(ctx *dgctx.DgContext, args ...any) {
	if dl.log.Core().Enabled(zap.ErrorLevel) {
		dl.withFields(ctx, nil, true).Error(fmt.Sprintln(args...))
	}
}

func (dl *DgLogger) Fatalln(ctx *dgctx.DgContext, args ...any) {
	dl.withFields(ctx, nil, true).Fatal(fmt.Sprintln(args...))
}

func (dl *DgLogger) Panicln(ctx *dgctx.DgContext, args ...any) {
	dl.withFields(ctx, nil, true).Panic(fmt.Sprintln(args...))
}

func (dl *DgLogger) Debugw(ctx *dgctx.DgContext, content string, fields map[string]any) {
	if dl.log.Core().Enabled(zap.DebugLevel) {
		dl.withFields(ctx, fields, false).Debug(content)
	}
}

func (dl *DgLogger) Infow(ctx *dgctx.DgContext, content string, fields map[string]any) {
	if dl.log.Core().Enabled(zap.InfoLevel) {
		dl.withFields(ctx, fields, false).Info(content)
	}
}

func (dl *DgLogger) Warnw(ctx *dgctx.DgContext, content string, fields map[string]any) {
	if dl.log.Core().Enabled(zap.WarnLevel) {
		dl.withFields(ctx, fields, false).Warn(content)
	}
}

func (dl *DgLogger) Errorw(ctx *dgctx.DgContext, content string, fields map[string]any) {
	if dl.log.Core().Enabled(zap.ErrorLevel) {
		dl.withFields(ctx, fields, true).Error(content)
	}
}

func (dl *DgLogger) Fatalw(ctx *dgctx.DgContext, content string, fields map[string]any) {
	dl.withFields(ctx, fields, true).Fatal(content)
}

func (dl *DgLogger) Panicw(ctx *dgctx.DgContext, content string, fields map[string]any) {
	dl.withFields(ctx, fields, true).Panic(content)
}

func SetExtraFields(ctx *dgctx.DgContext, fields map[string]any) {
	ctx.SetExtraKeyValue(extraFieldsKey, fields)
}

func (dl *DgLogger) withFields(ctx *dgctx.DgContext, fields map[string]any, printFileLine bool) *zap.Logger {
	allFields := []zap.Field{
		zap.String(constants.TraceId, ctx.TraceId),
	}

	if ctx.SpanId != "" {
		allFields = append(allFields, zap.String(constants.SpanId, ctx.SpanId))
	}

	if ctx.UserId > 0 {
		allFields = append(allFields, zap.Int64(constants.UID, ctx.UserId))
	}

	if len(fields) > 0 {
		for k, v := range fields {
			allFields = append(allFields, zap.Any(k, v))
		}
	}

	extraFields := ctx.GetExtraValue(extraFieldsKey)
	if extraFields != nil {
		fds := extraFields.(map[string]any)
		if len(fds) > 0 {
			for k, v := range fds {
				allFields = append(allFields, zap.Any(k, v))
			}
		}
	}

	if printFileLine {
		_, file, line, _ := runtime.Caller(3)
		allFields = append(allFields,
			zap.String("file", file),
			zap.String("line", strconv.Itoa(line)),
		)
	}

	return dl.log.With(allFields...)
}

func parseLevel(level string) zapcore.Level {
	switch level {
	case PanicLevel:
		return zap.PanicLevel
	case FatalLevel:
		return zap.FatalLevel
	case ErrorLevel:
		return zap.ErrorLevel
	case WarnLevel:
		return zap.WarnLevel
	case InfoLevel:
		return zap.InfoLevel
	case DebugLevel:
		return zap.DebugLevel
	default:
		return zap.DebugLevel
	}
}
