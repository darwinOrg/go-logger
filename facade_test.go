package dglogger

import (
	dgctx "github.com/darwinOrg/go-common/context"
	"testing"
)

func TestDebugf(t *testing.T) {
	Debugf(&dgctx.DgContext{TraceId: "123"}, "%s, %d", "abc", int64(789))
}

func TestInfof(t *testing.T) {
	Infof(&dgctx.DgContext{TraceId: "123"}, "%s, %d", "abc", int64(789))
}

func TestInfo(t *testing.T) {
	Info(&dgctx.DgContext{TraceId: "123"}, "abc", int64(789))
}

func TestInfoln(t *testing.T) {
	Infoln(&dgctx.DgContext{TraceId: "123"}, "abc", int64(789))
}

func TestLogEntry(t *testing.T) {
	ctx := &dgctx.DgContext{TraceId: "123"}
	Info(ctx, "erf", int64(789))
	Infoln(ctx, "abc", int64(456))
}
