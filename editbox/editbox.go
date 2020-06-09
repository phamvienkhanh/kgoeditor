package editbox

import (
	"io"
	"os"
	"time"

	"github.com/nsf/termbox-go"
)

type Editbox struct {
	currentX    int
	currentY    int
	width       int
	height      int
	currentCell termbox.Cell
	numLines    int
	fileName    string
	text        [][]rune
}

func (this *Editbox) Init() {
	this.text = make([][]rune, 1)
	this.currentX = 0
	this.currentY = 0
	this.numLines = 1
	this.width, this.height = termbox.Size()
	this.fileName = "./testest.txt"

}

func (this *Editbox) GetText() [][]rune {
	return this.text
}

func (this *Editbox) InsertTabAtCurrentPos() {
	for i := 0; i < 4; i++ {
		this.InsertCharAtCurrentPos(' ')
	}
}

func (this *Editbox) MoveCursorUp() {
	if this.currentY-1 >= 0 {
		this.currentY--
		if this.currentX >= len(this.text[this.currentY]) {
			this.currentX = len(this.text[this.currentY]) - 1
		}
	}
	this.updateCursor()
}

func (this *Editbox) MoveCursorDown() {
	if this.currentY+1 < this.numLines {
		this.currentY++
		if this.currentX >= len(this.text[this.currentY]) {
			this.currentX = len(this.text[this.currentY]) - 1
		}
	}
	this.updateCursor()
}

func (this *Editbox) MoveCursorLeft() {
	if this.currentX-1 >= 0 {
		this.currentX--
	} else {
		if this.currentY-1 >= 0 {
			this.currentY--
			this.currentX = len(this.text[this.currentY]) - 1
		}
	}
	this.updateCursor()
}

func (this *Editbox) MoveCursorRight() {
	if this.currentX+1 < len(this.text[this.currentY]) {
		this.currentX++
	} else {
		if this.currentY+1 < this.numLines {
			this.currentY++
			this.currentX = 0
		}
	}
	this.updateCursor()
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
	for _, line := range this.text {
		strBuffer += string(line)
	}
	this.writeToFile(this.fileName, strBuffer)
}

func (this *Editbox) BreakNewLine() {
	this.InsertCharAtCurrentPos('\n')
	this.numLines++

	newLine := this.text[this.currentY][this.currentX:]
	this.text[this.currentY] = this.text[this.currentY][this.currentX:]
	this.text = append(this.text, []rune{})
	this.currentY++

	for i := this.numLines; i < this.currentY; i-- {
		this.text[i] = this.text[i-1]
	}

	this.text[this.currentY] = newLine

	this.currentX = 0

	this.updateCursor()
	termbox.Flush()
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
	this.text[this.currentY] = append(this.text[this.currentY], 0)
	copy(this.text[this.currentY][this.currentX+1:], this.text[this.currentY][this.currentX:])
	this.text[this.currentY][this.currentX] = char
	this.currentX++
	this.updateCursor()
	termbox.Flush()
}

func (this *Editbox) AppendCharCurrentLine(char rune) {
	this.text[this.currentY] = append(this.text[this.currentY], char)
	this.currentX++
	this.updateCursor()
	termbox.Flush()
}

func (this *Editbox) ShowAllText() {
	for iLine, line := range this.text {
		for i, char := range line {
			termbox.SetCell(i, iLine, char, termbox.ColorDefault, termbox.ColorDefault)
		}
	}
	termbox.Flush()
}

func (this *Editbox) InsertAt(x, y int, char rune) {

}
