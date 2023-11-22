package dglogger

import (
	dgsys "github.com/darwinOrg/go-common/sys"
	"github.com/sirupsen/logrus"
)

func ConditionFatalf(condition bool, format string, args ...any) {
	if condition {
		logrus.Fatalf(format, args...)
	} else {
		logrus.Errorf(format, args...)
	}
}

func OnlyProdFatalf(format string, args ...any) {
	ConditionFatalf(dgsys.IsProd(), format, args...)
}
