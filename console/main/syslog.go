package main

import (
	"ghogo/console/shell/vterm"

	"github.com/sirupsen/logrus"
)

var logBuffer = make([]*logrus.Entry, 0)

type GeneralLogHook struct {
}

func (h *GeneralLogHook) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.TraceLevel,
		logrus.DebugLevel,
		logrus.InfoLevel,
		logrus.WarnLevel,
		logrus.ErrorLevel,
		logrus.PanicLevel,
		logrus.FatalLevel,
	}
}

func (h *GeneralLogHook) Fire(entry *logrus.Entry) error {
	logBuffer = append(logBuffer, entry)
	if vterm.Screen != nil {

		if display == DISPLAY_LOG {

			Repaint(vterm.Screen)
		}
	}
	return nil
}
