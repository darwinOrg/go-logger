package dglogger

import (
	"github.com/darwinOrg/go-common/constants"
	dgctx "github.com/darwinOrg/go-common/context"
	dgsys "github.com/darwinOrg/go-common/sys"
	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"io"
	"maps"
	"os"
	"runtime"
	"sort"
	"strconv"
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
	// TraceLevel level. Designates finer-grained informational events than the Debug.
	TraceLevel = "trace"
)

const (
	DefaultTimestampFormat = "2006-01-02 15:04:05.999999"
	DefaultFilename        = "app.log" // 日志文件路径
	DefaultMaxSize         = 100       // 每个日志文件的最大尺寸（MB）
	DefaultMaxBackups      = 10        // 保留旧日志文件的最大数量
	DefaultMaxAge          = 30        // 保留旧日志文件的最大天数
	DefaultCompress        = true      // 是否压缩/归档旧的日志文件
	extraFieldsKey         = "extraLogFields"
	logEntryKey            = "logEntry"
)

type DgLogger struct {
	log *logrus.Logger
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
	return &DgLogger{log: &logrus.Logger{
		Out: out,
		Formatter: &logrus.TextFormatter{
			DisableQuote:           true,
			FullTimestamp:          true,
			TimestampFormat:        timestampFormat,
			DisableSorting:         true,
			DisableLevelTruncation: true,
			PadLevelText:           false,
			SortingFunc: func(strings []string) {
				sort.Slice(strings, func(i, j int) bool {
					if strings[i] == "level" {
						return true
					}
					return false
				})
			},
		},
		Level: parseLevel(level),
	}}
}

func getDefaultLogLevel() string {
	if dgsys.IsProd() {
		return InfoLevel
	}

	return DebugLevel
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
	dl.withFields(ctx, nil, false).Debugf(format, args...)
}

func (dl *DgLogger) Infof(ctx *dgctx.DgContext, format string, args ...any) {
	dl.withFields(ctx, nil, false).Infof(format, args...)
}

func (dl *DgLogger) Warnf(ctx *dgctx.DgContext, format string, args ...any) {
	dl.withFields(ctx, nil, false).Warnf(format, args...)
}

func (dl *DgLogger) Errorf(ctx *dgctx.DgContext, format string, args ...any) {
	dl.withFields(ctx, nil, true).Errorf(format, args...)
}

func (dl *DgLogger) Fatalf(ctx *dgctx.DgContext, format string, args ...any) {
	dl.withFields(ctx, nil, true).Fatalf(format, args...)
}

func (dl *DgLogger) Panicf(ctx *dgctx.DgContext, format string, args ...any) {
	dl.withFields(ctx, nil, true).Panicf(format, args...)
}

func (dl *DgLogger) Debug(ctx *dgctx.DgContext, args ...any) {
	dl.withFields(ctx, nil, false).Debug(args...)
}

func (dl *DgLogger) Info(ctx *dgctx.DgContext, args ...any) {
	dl.withFields(ctx, nil, false).Info(args...)
}

func (dl *DgLogger) Warn(ctx *dgctx.DgContext, args ...any) {
	dl.withFields(ctx, nil, false).Warn(args...)
}

func (dl *DgLogger) Error(ctx *dgctx.DgContext, args ...any) {
	dl.withFields(ctx, nil, true).Error(args...)
}

func (dl *DgLogger) Fatal(ctx *dgctx.DgContext, args ...any) {
	dl.withFields(ctx, nil, true).Fatal(args...)
}

func (dl *DgLogger) Panic(ctx *dgctx.DgContext, args ...any) {
	dl.withFields(ctx, nil, true).Panic(args...)
}

func (dl *DgLogger) Debugln(ctx *dgctx.DgContext, args ...any) {
	dl.withFields(ctx, nil, false).Debugln(args...)
}

func (dl *DgLogger) Infoln(ctx *dgctx.DgContext, args ...any) {
	dl.withFields(ctx, nil, false).Infoln(args...)
}

func (dl *DgLogger) Warnln(ctx *dgctx.DgContext, args ...any) {
	dl.withFields(ctx, nil, false).Warnln(args...)
}

func (dl *DgLogger) Errorln(ctx *dgctx.DgContext, args ...any) {
	dl.withFields(ctx, nil, true).Errorln(args...)
}

func (dl *DgLogger) Fatalln(ctx *dgctx.DgContext, args ...any) {
	dl.withFields(ctx, nil, true).Fatalln(args...)
}

func (dl *DgLogger) Panicln(ctx *dgctx.DgContext, args ...any) {
	dl.withFields(ctx, nil, true).Panicln(args...)
}

func (dl *DgLogger) Debugw(ctx *dgctx.DgContext, content string, fields map[string]any) {
	dl.withFields(ctx, fields, false).Debug(content)
}

func (dl *DgLogger) Infow(ctx *dgctx.DgContext, content string, fields map[string]any) {
	dl.withFields(ctx, fields, false).Info(content)
}

func (dl *DgLogger) Warnw(ctx *dgctx.DgContext, content string, fields map[string]any) {
	dl.withFields(ctx, fields, false).Warn(content)
}

func (dl *DgLogger) Errorw(ctx *dgctx.DgContext, content string, fields map[string]any) {
	dl.withFields(ctx, fields, true).Error(content)
}

func (dl *DgLogger) Fatalw(ctx *dgctx.DgContext, content string, fields map[string]any) {
	dl.withFields(ctx, fields, true).Fatal(content)
}

func (dl *DgLogger) Panicw(ctx *dgctx.DgContext, content string, fields map[string]any) {
	dl.withFields(ctx, fields, true).Panic(content)
}

func (dl *DgLogger) SetLevel(level string) {
	dl.log.SetLevel(parseLevel(level))
}

func SetExtraFields(ctx *dgctx.DgContext, fields map[string]any) {
	ctx.SetExtraKeyValue(extraFieldsKey, fields)
}

func (dl *DgLogger) withFields(ctx *dgctx.DgContext, fields map[string]any, printFileLine bool) *log.Entry {
	if !printFileLine && len(fields) == 0 && ctx.GetExtraValue(logEntryKey) != nil {
		return ctx.GetExtraValue(logEntryKey).(*log.Entry)
	}

	allFields := log.Fields{constants.TraceId: ctx.TraceId}
	if ctx.UserId > 0 {
		allFields[constants.UID] = ctx.UserId
	}

	if len(fields) > 0 {
		maps.Copy(allFields, fields)
	}

	extraFields := ctx.GetExtraValue(extraFieldsKey)
	if extraFields != nil {
		fds := extraFields.(map[string]any)
		if len(fds) > 0 {
			maps.Copy(allFields, fds)
		}
	}

	if printFileLine {
		_, file, line, _ := runtime.Caller(3)
		allFields["file"] = file
		allFields["line"] = strconv.Itoa(line)

		return dl.log.WithFields(allFields)
	} else if len(fields) > 0 {
		return dl.log.WithFields(allFields)
	} else {
		entry := dl.log.WithFields(allFields)
		ctx.SetExtraKeyValue(logEntryKey, entry)
		return entry
	}
}

func parseLevel(level string) logrus.Level {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.DebugLevel
	}

	return logLevel
}
