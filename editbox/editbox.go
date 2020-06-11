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
	lines       *EditLine
	currentLine *EditLine
	fileName    string
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

func (this *Editbox) Init() {
	this.lines = &EditLine{0, []rune{}, nil, nil}
	this.currentLine = this.lines
	this.currentX = 0
	this.currentY = 0
	this.width, this.height = termbox.Size()
	this.fileName = "./testest.txt"

}

func (this *Editbox) InsertTabAtCurrentPos() {
	for i := 0; i < 4; i++ {
		this.InsertCharAtCurrentPos(' ')
	}
}

func (this *Editbox) MoveCursorUp() {
	if this.currentY-1 >= 0 {
		if this.currentLine.prevLine != nil {
			this.currentLine = this.currentLine.prevLine
		}
		this.currentY--
		if this.currentX >= this.currentLine.GetLen() {
			this.currentX = this.currentLine.GetLen()
		}
		this.updateCursor()
	}
}

func (this *Editbox) MoveCursorDown() {
	if this.currentY+1 < this.GetNumLines() {
		this.currentY++
		if this.currentLine.nextLine != nil {
			this.currentLine = this.currentLine.nextLine
		}
		if this.currentX >= this.currentLine.GetLen() {
			this.currentX = this.currentLine.GetLen()
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
			if this.currentLine.prevLine != nil {
				this.currentLine = this.currentLine.prevLine
			}
			this.currentX = this.currentLine.GetLen()
		}
	}
	this.updateCursor()
}

func (this *Editbox) MoveCursorRight() {
	if this.currentX < this.currentLine.GetLen() {
		this.currentX++
	} else {
		if this.currentY+1 < this.GetNumLines() {
			if this.currentLine.nextLine != nil {
				this.currentLine = this.currentLine.nextLine
			}
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
	iter := this.lines
	for iter != nil {
		strBuffer += string(iter.text) + "\n"
		iter = iter.nextLine
	}
	this.writeToFile(this.fileName, strBuffer)
}

func (this *Editbox) BreakNewLine() {
	if this.currentLine != nil {
		newLine := this.currentLine.SliceAt(this.currentX)

		iter := newLine.nextLine
		for iter != nil {
			iter.idLine += 1
			iter = iter.nextLine
		}
		if this.currentY+1 < this.height {
			this.currentY++
		}
		if this.currentLine.nextLine != nil {
			this.currentLine = this.currentLine.nextLine
		}
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
	if this.currentLine != nil {
		this.currentLine.InsertAt(this.currentX, char)
		this.currentX++
		this.updateCursor()
		termbox.Flush()
	}
}

func (this *Editbox) DeleteCharAtCurrentPos() {
	if this.currentLine != nil {
		var savePreLineLen = 0
		if this.currentY != 0 {
			savePreLineLen = this.currentLine.prevLine.GetLen()
		}
		this.currentLine.DeleteCharAt(this.currentX - 1)
		if this.currentX == 0 {
			if this.currentY != 0 {
				this.currentX = savePreLineLen
				this.currentY--
				if this.currentLine.prevLine != nil {
					this.currentLine = this.currentLine.prevLine
				}
				lastLine := this.GetLastLine()
				if lastLine != nil {
					for i := range lastLine.text {
						termbox.SetCell(i, lastLine.idLine+1, 0, termbox.ColorDefault, termbox.ColorDefault)
					}
				}
			}
		} else {
			this.currentX--
		}
	}
	this.updateCursor()
}

func (this *Editbox) GetLastLine() *EditLine {
	iter := this.lines
	if iter != nil && iter.nextLine != nil {
		iter = iter.nextLine
	}
	return iter
}

func (this *Editbox) ShowAllText() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.Flush()
	yPos := this.currentY
	iterLine := this.currentLine
	for yPos >= 0 {
		if iterLine != nil {
			for i, char := range iterLine.text {
				termbox.SetCell(i, yPos, char, termbox.ColorDefault, termbox.ColorDefault)
			}
			iterLine = iterLine.prevLine
		}
		yPos--
	}
	yPos = this.currentY + 1
	iterLine = this.currentLine
	for yPos < this.height {
		if iterLine != nil && iterLine.nextLine != nil {
			for i, char := range iterLine.text {
				termbox.SetCell(i, yPos, char, termbox.ColorDefault, termbox.ColorDefault)
			}
			iterLine = iterLine.nextLine
		}
		yPos++
	}
	termbox.Flush()
}
