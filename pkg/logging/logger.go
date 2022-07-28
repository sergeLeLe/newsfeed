package logging

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path"
	"runtime"
)

type logger struct {
	*logrus.Logger
}

func New(loggingLevel string, tsFormat string) (*logger, error) {
	l := logrus.New()

	l.Out = os.Stdout
	l.SetReportCaller(true)

	if lvl, err := logrus.ParseLevel(loggingLevel); err != nil {
		return nil, err
	} else {
		l.SetLevel(lvl)
	}

	l.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: tsFormat,
		PrettyPrint:     false,
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", path.Base(f.File), f.Line)
		},
	})
	return &logger{
		l,
	}, nil

}