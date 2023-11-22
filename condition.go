package dglogger

import (
	dgsys "github.com/darwinOrg/go-common/sys"
	"github.com/sirupsen/logrus"
)

func ConditionFatal(condition bool, args ...any) {
	if condition {
		logrus.Fatal(args...)
	} else {
		logrus.Error(args...)
	}
}

func ConditionFatalf(condition bool, format string, args ...any) {
	if condition {
		logrus.Fatalf(format, args...)
	} else {
		logrus.Errorf(format, args...)
	}
}

func OnlyProdFatal(args ...any) {
	ConditionFatal(dgsys.IsProd(), args...)
}

func QaOrProdFatal(args ...any) {
	ConditionFatal(dgsys.IsQa() || dgsys.IsProd(), args...)
}

func OnlyProdFatalf(format string, args ...any) {
	ConditionFatalf(dgsys.IsProd(), format, args...)
}

func QaOrProdFatalf(format string, args ...any) {
	ConditionFatalf(dgsys.IsQa() || dgsys.IsProd(), format, args...)
}
