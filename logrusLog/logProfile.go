package logrusLog

import (
	"github.com/siruspen/logrus"
)

var Log *logrus.Logger

func init() {
	Log = logrus.New()
}
