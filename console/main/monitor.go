package main

import (
	"ghogo/console/shell/vterm"
	"ghogo/util"
	"os"
	"strconv"

	"github.com/gdamore/tcell/v2"
	"github.com/sirupsen/logrus"
)

const (
	DISPLAY_LOG = iota
	DISPLAY_TERMINAL
)

var display = DISPLAY_LOG
var width, height int

const (
	PLAIN_DOWN_RIGHT         = "┌"
	PLAIN_DOWN_LEFT_RIGHT    = "┬"
	PLAIN_DOWN_LEFT          = "┐"
	PLAIN_UP_DOWN_RIGHT      = "├"
	PLAIN_UP_DOWN_LEFT_RIGHT = "┼"
	PLAIN_UP_DOWN_LEFT       = "┤"
	PLAIN_UP_RIGHT           = "└"
	PLAIN_UP_LEFT_RIGHT      = "┴"
	PLAIN_UP_LEFT            = "┘"
	PLAIN_LEFT_RIGHT         = "─"
)

var STYLE_RESET = tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
var STYLE_TRACE = tcell.StyleDefault.Background(tcell.ColorGray)
var STYLE_DEBUG = tcell.StyleDefault.Background(tcell.ColorDarkGreen)
var STYLE_INFO = tcell.StyleDefault.Background(tcell.ColorGreen)
var STYLE_WARN = tcell.StyleDefault.Background(tcell.ColorYellow)
var STYLE_ERROR = tcell.StyleDefault.Background(tcell.ColorRed)
var STYLE_PANIC = tcell.StyleDefault.Background(tcell.ColorDarkRed)
var STYLE_FATAL = tcell.StyleDefault.Background(tcell.ColorLightGray)

func VtermEvent(ev tcell.Event, s tcell.Screen) {

	switch ev := ev.(type) {
	case *tcell.EventResize:
		width, height = s.Size()
		s.Sync()
		Repaint(s)
	case *tcell.EventKey:
		if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {

			s.Fini()
			os.Exit(0)
		} else if ev.Modifiers() == 4 { //mod==4

			if ev.Rune() == 'l' || ev.Rune() == 'L' {
				display = DISPLAY_LOG
			} else if ev.Rune() == 't' || ev.Rune() == 'T' {
				display = DISPLAY_TERMINAL
			}
			Repaint(s)
		} else {

			mod, key, ch := ev.Modifiers(), ev.Key(), ev.Rune()

			vterm.DrawText(0, 0, 20, 10, tcell.StyleDefault, "                                       ")
			vterm.DrawText(0, 0, 20, 10, tcell.StyleDefault, "mod:"+strconv.FormatInt(int64(mod), 10)+" key:"+strconv.FormatInt(int64(key), 10)+" ch:"+string(ch))
		}
	}
}

func Repaint(s tcell.Screen) {
	s.Clear()
	for i := 0; i < width; i++ {
		vterm.DrawText(i, 1, width, 1, tcell.StyleDefault, "=")
	}
	vterm.DrawText(width-3, 0, width, 0, tcell.StyleDefault, strconv.Itoa(len(logBuffer)))
	if display == DISPLAY_TERMINAL {
		vterm.DrawText(0, 0, 35, 0, tcell.StyleDefault, "| >V[T]erminal |  System [L]og |")

	} else {
		vterm.DrawText(0, 0, 35, 0, tcell.StyleDefault, "|  V[T]erminal | >System [L]og |")
		RepaintLogWaterfall(s)
	}

	s.Sync()
}

func RepaintSubProcessList(s tcell.Screen) {

}

func RepaintLogWaterfall(s tcell.Screen) {
	if display != DISPLAY_LOG {
		return
	}

	length := len(logBuffer)
	idx := length - 1
	for y := height - 2; y >= 2; y-- {
		if idx < 0 {
			break
		}
		entry := (*logBuffer[idx])

		style := tcell.StyleDefault
		switch entry.Level {
		case logrus.TraceLevel:
			style = STYLE_TRACE
		case logrus.DebugLevel:
			style = STYLE_DEBUG
		case logrus.InfoLevel:
			style = STYLE_INFO
		case logrus.WarnLevel:
			style = STYLE_WARN
		case logrus.PanicLevel:
			style = STYLE_PANIC
		case logrus.FatalLevel:
			style = STYLE_FATAL

		}

		vterm.DrawText(0, y, width, y, tcell.StyleDefault, util.GetTimeStr(entry.Time))
		vterm.DrawText(20, y, 26, y, style, entry.Level.String())
		vterm.DrawText(26, y, width, y, STYLE_RESET, entry.Message)
		idx--
	}
}
