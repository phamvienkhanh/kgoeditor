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
	this.fileName = "./testest.txt"
}

func (this *Editbox) GetText() [][]rune {
	return this.text
}

func (this *Editbox) MoveCursorUp() {
	this.currentY--
	this.updateCursor()
}

func (this *Editbox) MoveCursorDown() {
	this.currentY++
	this.updateCursor()
}

func (this *Editbox) MoveCursorLeft() {
	this.currentX--
	this.updateCursor()
}

func (this *Editbox) MoveCursorRight() {
	this.currentX++
	this.updateCursor()
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

func (this *Editbox) AddNewLine() {
	this.AppendCharCurrentLine('\n')
	this.text = append(this.text, []rune{})
	this.numLines++
	this.currentX = 0
	this.currentY = this.numLines - 1
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
