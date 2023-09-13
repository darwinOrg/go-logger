package dglogger

import (
	"github.com/darwinOrg/go-common/constants"
	dgctx "github.com/darwinOrg/go-common/context"
	dgsys "github.com/darwinOrg/go-common/sys"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"sort"
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
	logEntryKey            = "logEntry"
)

type DgLogger struct {
	log *logrus.Logger
}

func DefaultDgLogger() *DgLogger {
	level := DebugLevel
	if dgsys.IsProd() {
		level = InfoLevel
	}
	return NewDgLogger(level, DefaultTimestampFormat, os.Stdout)
}

func NewDgLogger(level string, timestampFormat string, out io.Writer) *DgLogger {
	log := &logrus.Logger{
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
	}

	return &DgLogger{log: log}
}

func (dl *DgLogger) Debugf(ctx *dgctx.DgContext, format string, args ...any) {
	dl.withFields(ctx).Debugf(format, args...)
}

func (dl *DgLogger) Infof(ctx *dgctx.DgContext, format string, args ...any) {
	dl.withFields(ctx).Infof(format, args...)
}

func (dl *DgLogger) Warnf(ctx *dgctx.DgContext, format string, args ...any) {
	dl.withFields(ctx).Warnf(format, args...)
}

func (dl *DgLogger) Errorf(ctx *dgctx.DgContext, format string, args ...any) {
	dl.withFields(ctx).Errorf(format, args...)
}

func (dl *DgLogger) Fatalf(ctx *dgctx.DgContext, format string, args ...any) {
	dl.withFields(ctx).Fatalf(format, args...)
}

func (dl *DgLogger) Panicf(ctx *dgctx.DgContext, format string, args ...any) {
	dl.withFields(ctx).Panicf(format, args...)
}

func (dl *DgLogger) Debug(ctx *dgctx.DgContext, args ...any) {
	dl.withFields(ctx).Debug(args...)
}

func (dl *DgLogger) Info(ctx *dgctx.DgContext, args ...any) {
	dl.withFields(ctx).Info(args...)
}

func (dl *DgLogger) Warn(ctx *dgctx.DgContext, args ...any) {
	dl.withFields(ctx).Warn(args...)
}

func (dl *DgLogger) Error(ctx *dgctx.DgContext, args ...any) {
	dl.withFields(ctx).Error(args...)
}

func (dl *DgLogger) Fatal(ctx *dgctx.DgContext, args ...any) {
	dl.withFields(ctx).Fatal(args...)
}

func (dl *DgLogger) Panic(ctx *dgctx.DgContext, args ...any) {
	dl.withFields(ctx).Panic(args...)
}

func (dl *DgLogger) Debugln(ctx *dgctx.DgContext, args ...any) {
	dl.withFields(ctx).Debugln(args...)
}

func (dl *DgLogger) Infoln(ctx *dgctx.DgContext, args ...any) {
	dl.withFields(ctx).Infoln(args...)
}

func (dl *DgLogger) Warnln(ctx *dgctx.DgContext, args ...any) {
	dl.withFields(ctx).Warnln(args...)
}

func (dl *DgLogger) Errorln(ctx *dgctx.DgContext, args ...any) {
	dl.withFields(ctx).Errorln(args...)
}

func (dl *DgLogger) Fatalln(ctx *dgctx.DgContext, args ...any) {
	dl.withFields(ctx).Fatalln(args...)
}

func (dl *DgLogger) Panicln(ctx *dgctx.DgContext, args ...any) {
	dl.withFields(ctx).Panicln(args...)
}

func (dl *DgLogger) SetLevel(level string) {
	dl.log.SetLevel(parseLevel(level))
}

func (dl *DgLogger) withFields(ctx *dgctx.DgContext) *log.Entry {
	if ctx.GetExtraValue(logEntryKey) != nil {
		return ctx.GetExtraValue(logEntryKey).(*log.Entry)
	}

	if ctx.GoId == 0 {
		ctx.GoId = dgsys.QuickGetGoRoutineId()
	}

	entry := dl.log.WithFields(log.Fields{
		constants.GoId:    ctx.GoId,
		constants.TraceId: ctx.TraceId,
	})

	ctx.SetExtraKeyValue(logEntryKey, entry)

	return entry
}

func parseLevel(level string) logrus.Level {
	logLevel, err := logrus.ParseLevel(level)
	if err != nil {
		logLevel = logrus.DebugLevel
	}

	return logLevel
}
