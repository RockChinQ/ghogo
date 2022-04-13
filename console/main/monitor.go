package main

import (
	"ghogo/console/shell/vterm"
	"ghogo/util"
	"ghogo/util/puppet"
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

var prompt = "console"

var cursorX = 0

var buffer = make([]rune, 0)

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

var STYLE_PROMPT = tcell.StyleDefault.Background(tcell.ColorLightGreen)

var STYLE_CURSOR = tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack)

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

			vterm.DrawText(32, 0, 55, 0, tcell.StyleDefault, []rune("                                       "), false)
			vterm.DrawText(32, 0, 55, 0, tcell.StyleDefault, []rune(" mod:"+strconv.FormatInt(int64(mod), 10)+" key:"+strconv.FormatInt(int64(key), 10)+" ch:"+string(ch)+"|"), false)

			if display == DISPLAY_TERMINAL {

				if key == 256 { //insert symbol
					InsertBuffer(s, ch)
				} else if mod == 0 && key == 8 { //backspace
					Backspace(s)
				} else if mod == 0 && key == 13 { //enter
					logrus.Trace("Enter:", GetBufferText())
					buffer = make([]rune, 0)
					cursorX = 0
					RepaintBuffer(s)
				} else if mod == 0 && key == 260 { //left
					MoveCursor(s, -1)
				} else if mod == 0 && key == 259 { //right
					MoveCursor(s, 1)
				}
				RepaintDebug(s)
				s.Sync()
			}
		}
	}
}

func Repaint(s tcell.Screen) {
	s.Clear()
	for i := 0; i < width; i++ {
		vterm.DrawText(i, 1, width, 1, tcell.StyleDefault, []rune("="), false)
	}
	if display == DISPLAY_TERMINAL {
		vterm.DrawText(0, 0, 35, 0, tcell.StyleDefault, []rune("| >V[T]erminal |  System [L]og |"), false)
		RepaintVTerminalPanel(s)
	} else {
		vterm.DrawText(0, 0, 35, 0, tcell.StyleDefault, []rune("|  V[T]erminal | >System [L]og |"), false)
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

		vterm.DrawText(0, y, width, y, tcell.StyleDefault, []rune(util.GetTimeStr(entry.Time)), false)
		vterm.DrawText(20, y, 26, y, style, []rune(entry.Level.String()), false)
		vterm.DrawText(26, y, width, y, STYLE_RESET, []rune(entry.Message), false)
		idx--
	}
}

//current focused subprocess,any input through vterm will be redirected to this
var focusedSubprocess puppet.SubProcess

func RepaintVTerminalPanel(s tcell.Screen) {
	RepaintBuffer(s)
}

func RepaintBuffer(s tcell.Screen) {

	for i := 0; i < width; i++ {
		vterm.DrawText(i, height-1, width, height-1, tcell.StyleDefault, []rune(" "), false)
	}

	//put prompt
	vterm.DrawText(0, height-1, len(prompt)+6, height-1, STYLE_PROMPT, []rune("GHO "+prompt+" >"), false)

	//buffer
	vterm.DrawText(len(prompt)+6+1, height-1, width, height-1, tcell.StyleDefault, buffer, false)

	//cursor
	DrawCursor(s)

	s.Sync()
}

func MoveCursor(s tcell.Screen, delta int) {
	ClearCursor(s)
	if delta == -1 && cursorX > 0 {
		if buffer[cursorX-1] > 255 {
			cursorX--
		}
		cursorX -= 1
	} else if delta == 1 && cursorX < len(buffer) {
		if buffer[cursorX] > 255 {
			cursorX++
		}
		cursorX += 1
	}
	//cursor
	DrawCursor(s)
}

func InsertBuffer(s tcell.Screen, ch rune) {

	logrus.Trace("Insert rune:", ch)

	if ch <= 255 {

		buffer = append(buffer[:cursorX], append([]rune{ch}, buffer[cursorX:]...)...)
		cursorX += 1
		//绘制insert位置之后的所有字符
		vterm.DrawText(len(prompt)+6+cursorX, height-1, len(prompt)+6+len(buffer), height-1, tcell.StyleDefault, buffer[cursorX-1:], true)
	} else if ch > 255 {
		buffer = append(buffer[:cursorX], append([]rune{ch}, buffer[cursorX:]...)...)
		//绘制insert位置之后的所有字符
		cursorX += 1
		vterm.DrawText(len(prompt)+6+cursorX, height-1, len(prompt)+6+len(buffer), height-1, tcell.StyleDefault, buffer[cursorX-1:], true)
		buffer = append(buffer[:cursorX], append([]rune{ch}, buffer[cursorX:]...)...)
		vterm.DrawText(len(prompt)+6+cursorX, height-1, len(prompt)+6+len(buffer), height-1, tcell.StyleDefault, buffer[cursorX-1:], true)
		cursorX += 1
	}

	//cursor
	DrawCursor(s)
}

func Backspace(s tcell.Screen) {
	if cursorX != 0 {
		ClearCursor(s)
		if buffer[cursorX-1] <= 255 {

			//delete end
			vterm.DrawText(len(prompt)+6+len(buffer), height-1, len(prompt)+6+len(buffer)-1, height-1, tcell.StyleDefault, []rune(" "), false)

			buffer = append(buffer[:cursorX-1], buffer[cursorX:]...)

			//draw text after cursor
			cursorX--
		} else if buffer[cursorX-1] > 255 {
			//delete end
			vterm.DrawText(len(prompt)+6+len(buffer), height-1, len(prompt)+6+len(buffer)-1, height-1, tcell.StyleDefault, []rune(" "), false)

			buffer = append(buffer[:cursorX-1], buffer[cursorX:]...)

			//draw text after cursor
			cursorX--

			//delete end
			vterm.DrawText(len(prompt)+6+len(buffer), height-1, len(prompt)+6+len(buffer)-1, height-1, tcell.StyleDefault, []rune(" "), false)

			buffer = append(buffer[:cursorX-1], buffer[cursorX:]...)

			//draw text after cursor
			cursorX--
		}

		//cursor
		DrawCursor(s)
	}

}

func DrawCursor(s tcell.Screen) {
	if cursorX >= len(buffer) {
		vterm.DrawText(len(prompt)+7+len(buffer), height-1, len(prompt)+7+len(buffer), height-1, STYLE_CURSOR, []rune(" "), false)
	} else {
		vterm.DrawText(len(prompt)+7+cursorX, height-1, len(prompt)+7+cursorX, height-1, STYLE_CURSOR, buffer[cursorX:cursorX+1], false)
	}
}
func ClearCursor(s tcell.Screen) {
	if cursorX >= len(buffer) {
		vterm.DrawText(len(prompt)+7+len(buffer), height-1, len(prompt)+7+len(buffer), height-1, tcell.StyleDefault, []rune(" "), false)
	} else {
		vterm.DrawText(len(prompt)+7+cursorX, height-1, len(prompt)+7+cursorX, height-1, tcell.StyleDefault, buffer[cursorX:cursorX+1], false)
	}
}

func RepaintDebug(s tcell.Screen) {

	vterm.DrawText(width-3, 0, width, 0, tcell.StyleDefault, []rune(strconv.Itoa(len(buffer))), false)
	//cursor
	DrawCursor(s)
}

func GetBufferText() string {
	text := make([]rune, 0)

	for i := 0; i < len(buffer); i++ {
		if buffer[i] <= 255 {
			text = append(text, buffer[i])
		} else if buffer[i] > 255 {
			text = append(text, buffer[i])
			i++
		}
	}
	return string(text)
}
