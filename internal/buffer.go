package internal

import (
	"os"
	"strconv"
	"strings"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Modes
const (
	Insert = iota
	Normal
)

// Represents the history of the file, meaning the lines changed and it's index.
type Change struct {
	Lines []string
	Index int
}
type Cursor struct {
	X       int
	Y       int
	XOffset int
	YOffset int
}
type Buffer struct {
	Path         string
	Mode         int
	Lines        []string
	BufferCursor Cursor
	UndoList     []Change
	Padding      int
}

func absInt(num int) int {
	if num < 0 {
		return num * -1
	}
	return num
}

func NewBuffer(path string) (Buffer, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return Buffer{}, err
	}
	cursor := Cursor{X: 0, Y: 0, YOffset: 0, XOffset: 0}
	undos := make([]Change, 0)
	lines := strings.Split(string(content), "\n")
	return Buffer{Path: path, Mode: Insert, Lines: lines, BufferCursor: cursor, UndoList: undos, Padding: 10}, nil
}

func (b *Buffer) RenderBuffer(fontSize int, fontColor, highlight rl.Color, font rl.Font) error {
	for i, line := range b.Lines {
		var lineNumber string
		lineNumberColor := rl.White
		if i != b.BufferCursor.Y {
			lineNumber = strconv.Itoa(absInt(i - b.BufferCursor.Y))
		} else {
			lineNumber = strconv.Itoa(i + 1)
			lineNumberColor = highlight
		}
		position := rl.Vector2{X: float32(b.Padding + 1), Y: float32(i * fontSize)}
		rl.DrawTextEx(font, lineNumber, position, float32(fontSize), float32(fontSize), lineNumberColor)
		position.X = float32(b.Padding + fontSize)
		rl.DrawTextEx(font, line, position, float32(fontSize), 0, fontColor)
	}
	return nil
}
func (b *Buffer) ListenInput() {
	if b.Mode == Insert {
		if rl.IsKeyPressed(rl.KeyJ) && b.BufferCursor.Y < len(b.Lines)-1 {
			b.BufferCursor.Y++
		}
		if rl.IsKeyPressed(rl.KeyK) && b.BufferCursor.Y > 0 {
			b.BufferCursor.Y--
		}
	}
}
