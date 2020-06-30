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
	offsetLeft  int
	offsetTop   int
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
	this.offsetLeft = 3
	this.offsetTop = 2
	this.currentX = this.offsetLeft
	this.currentY = this.offsetTop
	this.width, this.height = termbox.Size()
	this.fileName = "./testest.txt"
	termbox.SetInputMode(termbox.InputAlt | termbox.InputMouse | termbox.InputEsc)
	termbox.SetCursor(this.offsetLeft, this.offsetTop)
	termbox.Flush()
}

func (this *Editbox) InsertTabAtCurrentPos() {
	for i := 0; i < 4; i++ {
		this.InsertCharAtCurrentPos(' ')
	}
}

func (this *Editbox) MoveCursorUp() {
	if this.currentY-1-this.offsetTop >= 0 {
		this.currentY--
	}

	if this.currentLine.prevLine != nil {
		this.currentLine = this.currentLine.prevLine
		this.ShowAllText()
	}
	if this.currentX-this.offsetLeft >= this.currentLine.GetLen() {
		this.currentX = this.currentLine.GetLen() + this.offsetLeft
	}
	this.updateCursor()
}

func (this *Editbox) MoveCursorDown() {
	if this.currentY+1 <= this.height && this.currentLine.nextLine != nil {
		this.currentY++
	}
	if this.currentLine.nextLine != nil {
		this.currentLine = this.currentLine.nextLine
		this.ShowAllText()
	}
	if this.currentX-this.offsetLeft >= this.currentLine.GetLen() {
		this.currentX = this.currentLine.GetLen() + this.offsetLeft
	}
	this.updateCursor()
}

func (this *Editbox) MoveCursorLeft() {
	if this.currentX-1-this.offsetLeft >= 0 {
		this.currentX--
	} else {
		if this.currentY-1-this.offsetTop >= 0 {
			if this.currentY-1-this.offsetTop >= 0 {
				this.currentY--
			}
			if this.currentLine.prevLine != nil {
				this.currentLine = this.currentLine.prevLine
			}
			this.currentX = this.currentLine.GetLen() + this.offsetLeft
		}
	}
	this.updateCursor()
}

func (this *Editbox) MoveCursorRight() {
	if this.currentX-this.offsetLeft < this.currentLine.GetLen() {
		this.currentX++
	} else {
		if this.currentY+1-this.offsetTop < this.GetNumLines() {
			if this.currentLine.nextLine != nil {
				this.currentLine = this.currentLine.nextLine
			}
			if this.currentY+1 < this.height {
				this.currentY++
			}
			this.currentX = this.offsetLeft
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
		newLine := this.currentLine.SliceAt(this.currentX - this.offsetLeft)

		iter := newLine.nextLine
		for iter != nil {
			iter.idLine += 1
			iter = iter.nextLine
		}
		if this.currentY+1 < this.height {
			this.currentY++
		}
		this.currentLine = newLine
		this.currentX = this.offsetLeft
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
		this.currentLine.InsertAt(this.currentX-this.offsetLeft, char)
		this.currentX++
		this.updateCursor()
		termbox.Flush()
	}
}

func (this *Editbox) DeleteCharAtCurrentPos() {
	if this.currentLine != nil {
		var savePreLineLen = 0
		if this.currentY != this.offsetTop {
			savePreLineLen = this.currentLine.prevLine.GetLen()
		}
		this.currentLine.DeleteCharAt(this.currentX - 1 - this.offsetLeft)
		if this.currentX <= this.offsetLeft {
			if this.currentY != this.offsetTop {
				this.currentX = savePreLineLen + this.offsetLeft
				if this.currentY-1-this.offsetTop >= 0 {
					this.currentY--
				}
				if this.currentLine.prevLine != nil {
					this.currentLine = this.currentLine.prevLine
				}
				lastLine := this.GetLastLine()
				if lastLine != nil {
					for i := range lastLine.text {
						termbox.SetCell(i+this.offsetLeft, lastLine.idLine+1, 0, termbox.ColorDefault, termbox.ColorDefault)
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

func countNums(num int) int {
	count := 1

	for num >= 10 {
		num /= 10
		count++
	}

	return count
}

func (this *Editbox) ShowNumLine(yPos int, line *EditLine) {
	lineId := line.GetLineId()
	lenLineId := len(lineId)

	newOffset := countNums(this.GetNumLines()) + 3

	// print space before num line
	numSpace := newOffset - lenLineId - 2

	// print num line
	for i, num := range lineId {
		termbox.SetCell(i+numSpace, yPos, num, termbox.ColorYellow, termbox.ColorDefault)
	}

	termbox.SetCell(lenLineId+numSpace, yPos, 9475, termbox.ColorMagenta, termbox.ColorBlack)
	termbox.SetCell(lenLineId+numSpace+1, yPos, ' ', termbox.ColorDefault, termbox.ColorDefault)

	if this.offsetLeft != newOffset {
		this.currentX += newOffset - this.offsetLeft
	}

	this.offsetLeft = newOffset
}

func (this *Editbox) ShowAllText() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.Flush()
	yPos := this.currentY
	iterLine := this.currentLine
	for yPos >= this.offsetTop {
		if iterLine != nil {
			this.ShowNumLine(yPos, iterLine)
			for i, char := range iterLine.text {
				termbox.SetCell(i+this.offsetLeft, yPos, char, termbox.ColorDefault, termbox.ColorDefault)
			}
			iterLine = iterLine.prevLine
		}
		yPos--
	}
	yPos = this.currentY
	iterLine = this.currentLine
	for yPos < this.height {
		if iterLine != nil {
			this.ShowNumLine(yPos, iterLine)
			for i, char := range iterLine.text {
				termbox.SetCell(i+this.offsetLeft, yPos, char, termbox.ColorDefault, termbox.ColorDefault)
			}
			iterLine = iterLine.nextLine
		}
		yPos++
	}
	termbox.Flush()
}
