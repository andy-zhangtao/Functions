package flogs

import "github.com/sirupsen/logrus"

func init() {
	logrus.SetLevel(logrus.InfoLevel)
}

func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}
