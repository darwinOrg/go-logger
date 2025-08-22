package dglogger

import (
	"errors"
	"testing"

	dgctx "github.com/darwinOrg/go-common/context"
)

func TestDebugf(t *testing.T) {
	ctx := dgctx.SimpleDgContext()
	SetExtraFields(ctx, map[string]any{"key": "value"})
	Debugf(ctx, "%s, %d", "abc", int64(789))
}

func TestInfof(t *testing.T) {
	Infof(dgctx.SimpleDgContext(), "%s, %d", "abc", int64(789))
}

func TestInfo(t *testing.T) {
	Info(dgctx.SimpleDgContext(), "abc", int64(789))
}

func TestInfoln(t *testing.T) {
	Infoln(dgctx.SimpleDgContext(), "abc", int64(789))
}

func TestErrorf(t *testing.T) {
	Errorf(dgctx.SimpleDgContext(), "abc: %d", int64(789))
}

func TestLogEntry(t *testing.T) {
	ctx := dgctx.SimpleDgContext()
	Info(ctx, "erf", int64(789))
	Infoln(ctx, "abc", int64(456))
	Infow(ctx, "hij", "key1", 1, "key2", "2", "err", errors.New("illegal"))
}
