package editbox

import (
	"io"
	"os"
	"time"

	"github.com/nsf/termbox-go"
)

type Editbox struct {
	currentX int
	currentY int
	width    int
	height   int
	lines    *EditLine
	fileName string
}

func (this *Editbox) GetNumLines() int {
	countLines := 0
	iter := this.lines
	for iter != nil {
		countLines++
		iter = iter.nextLine
	}
	return countLines
}

func (this *Editbox) GetCurrentLine() *EditLine {
	return this.GetLineAt(this.currentY)
}

func (this *Editbox) Init() {
	this.lines = &EditLine{0, []rune{}, nil, nil}
	this.currentX = 0
	this.currentY = 0
	this.width, this.height = termbox.Size()
	this.fileName = "./testest.txt"

}

func (this *Editbox) GetLenLine(index int) int {
	iter := this.lines
	for iter != nil {
		if iter.idLine == index {
			return len(iter.text)
		}
		iter = iter.nextLine
	}

	return -1
}

func (this *Editbox) GetLineAt(index int) *EditLine {
	iter := this.lines
	for iter != nil {
		if iter.idLine == index {
			return iter
		}
		iter = iter.nextLine
	}

	return nil
}

func (this *Editbox) InsertTabAtCurrentPos() {
	for i := 0; i < 4; i++ {
		this.InsertCharAtCurrentPos(' ')
	}
}

func (this *Editbox) MoveCursorUp() {
	if this.currentY-1 >= 0 {
		this.currentY--
		if this.currentX >= this.GetLenLine(this.currentY) {
			this.currentX = this.GetLenLine(this.currentY) - 1
		}
		this.updateCursor()
	}
}

func (this *Editbox) MoveCursorDown() {
	if this.currentY+1 < this.GetNumLines() {
		this.currentY++
		if this.currentX >= this.GetLenLine(this.currentY) {
			this.currentX = this.GetLenLine(this.currentY) - 1
		}
		this.updateCursor()
	}
}

func (this *Editbox) MoveCursorLeft() {
	if this.currentX-1 >= 0 {
		this.currentX--
	} else {
		if this.currentY-1 >= 0 {
			this.currentY--
			this.currentX = this.GetLenLine(this.currentY) - 1
		}
	}
	this.updateCursor()
}

func (this *Editbox) MoveCursorRight() {
	if this.GetLenLine(this.currentY) != 0 {
		if this.currentX < this.GetLenLine(this.currentY) {
			this.currentX++
		} else {
			if this.currentY+1 < this.GetNumLines() {
				this.currentY++
				this.currentX = 0
			}
		}
		this.updateCursor()
	}
}

func (this *Editbox) SetWindowSize() {
	this.width, this.height = termbox.Size()
}

func (this *Editbox) writeToFile(filename string, data string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.WriteString(file, data)
	if err != nil {
		return err
	}
	return file.Sync()
}

func (this *Editbox) SaveToFile() {
	strBuffer := ""
	iter := this.lines
	for iter != nil {
		strBuffer += string(iter.text) + "\n"
		iter = iter.nextLine
	}
	this.writeToFile(this.fileName, strBuffer)
}

func (this *Editbox) BreakNewLine() {
	currentLine := this.GetCurrentLine()
	if currentLine != nil {
		newLine := currentLine.SliceAt(this.currentX)

		iter := newLine.nextLine
		for iter != nil {
			iter.idLine += 1
			iter = iter.nextLine
		}
		this.currentY++
		this.currentX = 0
		this.updateCursor()
		termbox.Flush()
	}

}

func (this *Editbox) BlindCursor() {
	for {
		termbox.SetCursor(this.currentX, this.currentY)
		termbox.Flush()
		time.Sleep(500 * time.Millisecond)
		termbox.HideCursor()
		termbox.Flush()
		time.Sleep(500 * time.Millisecond)
	}
}

func (this *Editbox) updateCursor() {
	termbox.SetCursor(this.currentX, this.currentY)
	termbox.Flush()
}

func (this *Editbox) InsertCharAtCurrentPos(char rune) {
	currentLine := this.GetLineAt(this.currentY)
	if currentLine != nil {
		currentLine.InsertAt(this.currentX, char)
		this.currentX++
		this.updateCursor()
		termbox.Flush()
	}
}

func (this *Editbox) DeleteCharAtCurrentPos() {
	currentLine := this.GetCurrentLine()
	if currentLine != nil {
		currentLine.DeleteCharAt(this.currentX)
		if this.currentX == 0 {
			if this.currentY != 0 {
				this.currentX = currentLine.prevLine.GetLen()
				this.currentY--
			}
		} else {
			this.currentX--
		}
	}
	this.updateCursor()
}

func (this *Editbox) ShowAllText() {
	iter := this.lines
	for iter != nil {
		for i, char := range iter.text {
			termbox.SetCell(i, iter.idLine, char, termbox.ColorDefault, termbox.ColorDefault)
		}
		iter = iter.nextLine
	}
	termbox.Flush()
}

func (this *Editbox) InsertAt(x, y int, char rune) {

}
