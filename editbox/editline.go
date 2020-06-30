package editbox

import (
	"strconv"

	"github.com/nsf/termbox-go"
)

type EditLine struct {
	idLine   int
	text     []rune
	prevLine *EditLine
	nextLine *EditLine
}

func (this *EditLine) GetLen() int {
	return len(this.text)
}

func (this *EditLine) GetLineId() []rune {
	strId := strconv.FormatInt(int64(this.idLine), 10)
	return []rune(strId)
}

func (this *EditLine) InsertAt(index int, char rune) {
	if this.GetLen() == 0 || index == this.GetLen() {
		this.text = append(this.text, 0)
		this.text[index] = char
	} else {
		if index < this.GetLen() {
			this.text = append(this.text, 0)
			copy(this.text[index+1:], this.text[index:])
			this.text[index] = char
		}
	}
}

func (this *EditLine) DeleteLine() {
	preLine := this.prevLine
	nexLine := this.nextLine
	if preLine != nil {
		preLine.nextLine = nexLine
	}
	if nexLine != nil {
		nexLine.prevLine = preLine

		// reset index
		for nexLine != nil {
			nexLine.idLine--
			nexLine = nexLine.nextLine
		}
	}
	this = nil
}

func (this *EditLine) DeleteCharAt(index int) {
	lenLine := this.GetLen()
	if index < 0 {
		preLine := this.prevLine
		if preLine != nil {
			preLine.text = append(preLine.text, this.text...)
			this.DeleteLine()
		}
	} else if index < lenLine {
		copy(this.text[index:], this.text[index+1:])
		this.text[lenLine-1] = 0
		this.text = this.text[:lenLine-1]
		termbox.SetCell(lenLine-1, this.idLine, 0, termbox.ColorDefault, termbox.ColorDefault)
	} else {
		this.text = this.text[:index-1]
		termbox.SetCell(lenLine-1, this.idLine, 0, termbox.ColorDefault, termbox.ColorDefault)
	}
}

func (this *EditLine) SliceAt(index int) *EditLine {
	lenLine := this.GetLen()
	var newLine *EditLine
	var newText []rune
	if index < lenLine {
		// Clean up the current line on the screen
		for i := index; i < lenLine; i++ {
			termbox.SetCell(i, this.idLine, 0, termbox.ColorDefault, termbox.ColorDefault)
		}
		newText = make([]rune, lenLine-index)
		copy(newText, this.text[index:])
		this.text = this.text[:index]
	} else {
		newText = make([]rune, 0)
	}

	// Clean up the next line on the screen where the new line will be added
	if this.nextLine != nil {
		nextLineLen := len(this.nextLine.text)
		for i := 0; i < nextLineLen; i++ {
			termbox.SetCell(i, this.nextLine.idLine, 0, termbox.ColorDefault, termbox.ColorDefault)
		}
	}

	newLine = &EditLine{this.idLine + 1, newText, this, this.nextLine}
	if this.nextLine != nil && this.nextLine.prevLine != nil {
		this.nextLine.prevLine = newLine
	}
	this.nextLine = newLine

	return newLine
}
