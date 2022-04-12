package vterm

import (
	"io"

	"github.com/gdamore/tcell/v2"
	"github.com/sirupsen/logrus"
)

var Screen tcell.Screen

func LaunchVterm(hook func(tcell.Event, tcell.Screen)) {
	s, err := tcell.NewScreen()
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"location": "shell/vterm.go",
		}).Fatal("failed to launch vterm:" + err.Error())
	}
	if err := s.Init(); err != nil {
		logrus.WithFields(logrus.Fields{
			"location": "shell/vterm.go",
		}).Fatal("failed to launch vterm:" + err.Error())
	}

	Screen = s

	// Set default text style
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)

	// Clear screen
	s.Clear()

	logrus.SetOutput(io.Discard)

	for {
		// Update screen
		s.Show()

		// Poll event
		ev := s.PollEvent()

		// Process event
		hook(ev, Screen)
	}
}

func DrawText(x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		Screen.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}
