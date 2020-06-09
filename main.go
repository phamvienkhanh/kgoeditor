package main

import (
	"goeditor/editbox"

	"github.com/nsf/termbox-go"
)

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputAlt | termbox.InputMouse | termbox.InputEsc)
	termbox.SetCursor(0, 0)
	termbox.Flush()

	curEdit := editbox.Editbox{}
	curEdit.Init()

	go curEdit.BlindCursor()

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowRight:
				curEdit.MoveCursorRight()

			case termbox.KeyArrowLeft:
				curEdit.MoveCursorLeft()

			case termbox.KeyArrowDown:
				curEdit.MoveCursorDown()

			case termbox.KeyArrowUp:
				curEdit.MoveCursorUp()

			case termbox.KeyEnter:
				curEdit.BreakNewLine()

			case termbox.KeyTab:
				curEdit.InsertTabAtCurrentPos()

			case termbox.KeyEsc:
				curEdit.SaveToFile()
				break mainloop

			case termbox.KeySpace:
				curEdit.InsertCharAtCurrentPos(' ')
				curEdit.ShowAllText()

			default:
				curEdit.InsertCharAtCurrentPos(ev.Ch)
				curEdit.ShowAllText()
			}
		case termbox.EventInterrupt:
			break mainloop

		case termbox.EventResize:
			curEdit.SetWindowSize()

		case termbox.EventError:
			panic(ev.Err)
		}
	}
}
