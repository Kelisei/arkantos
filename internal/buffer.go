package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode/utf8"

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

func (c *Cursor) DrawCursor(fontSize, padding int, color rl.Color, b Buffer, font rl.Font) {
	lineStart := padding + fontSize
	line := b.Lines[c.Y]

	xStart := int32(lineStart + int(rl.MeasureTextEx(font, line[:c.X], float32(fontSize), 0).X))
	yStart := int32(fontSize * c.Y)
	xEnd := xStart
	yEnd := yStart + int32(fontSize)

	rl.DrawLine(xStart, yStart, xEnd, yEnd, color)
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
	fmt.Println(path)
	path, err := filepath.Abs(path)
	if err != nil {
		return Buffer{}, err
	}
	fmt.Println(path)
	var content []byte
	content, err = os.ReadFile(path)
	if err != nil {
		return Buffer{}, err
	}
	cursor := Cursor{X: 0, Y: 0, YOffset: 0, XOffset: 0}
	undos := make([]Change, 0)
	lines := strings.Split(string(content), "\n")
	return Buffer{Path: path, Mode: Insert, Lines: lines, BufferCursor: cursor, UndoList: undos, Padding: 10}, nil
}

func (b *Buffer) RenderBuffer(fontSize int, fontColor, highlight rl.Color, font rl.Font) error {
	b.BufferCursor.DrawCursor(fontSize, b.Padding, highlight, *b, font)
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
		rl.DrawTextEx(font, lineNumber, position, float32(fontSize), 0, lineNumberColor)
		position.X = float32(b.Padding + fontSize)
		rl.DrawTextEx(font, line, position, float32(fontSize), 0, fontColor)
	}
	drawBottomBar(fontSize, b, font, highlight)
	return nil
}

// Given certain parameters and the size of the screen, draws the current buffer,
// the cursor position and the current mode.
func drawBottomBar(fontSize int, b *Buffer, font rl.Font, highlight rl.Color) {
	height := rl.GetScreenHeight()
	width := rl.GetScreenWidth()
	bottomPos := float32(height - fontSize)

	pathPos := rl.Vector2{X: float32(b.Padding), Y: bottomPos}
	pathStr := "Buffer :" + b.Path
	rl.DrawTextEx(font, pathStr, pathPos, float32(fontSize), 0, highlight)

	currentMode := ""
	switch b.Mode {
	case Insert:
		currentMode = "Insert"
	case Normal:
		currentMode = "Normal"
	default:
		currentMode = "Unknown"
	}
	modePos := rl.Vector2{X: float32(width - utf8.RuneCountInString(currentMode)*fontSize), Y: bottomPos}
	rl.DrawTextEx(font, "--"+currentMode+"--", modePos, float32(fontSize), 0, highlight)

	cursorPosStr := strconv.Itoa(b.BufferCursor.X+1) + "," + strconv.Itoa(b.BufferCursor.Y+1)
	cursorPos := rl.Vector2{X: float32(rl.MeasureText(pathStr, int32(fontSize))), Y: bottomPos}
	rl.DrawTextEx(font, cursorPosStr, cursorPos, float32(fontSize), 0, highlight)
}

func wrapCursor(previousLineLen int, b *Buffer) {
	if b.BufferCursor.X > utf8.RuneCountInString(b.Lines[b.BufferCursor.Y]) || b.BufferCursor.X == previousLineLen {
		b.BufferCursor.X = utf8.RuneCountInString(b.Lines[b.BufferCursor.Y])
	}
}

func (b *Buffer) ListenInput() {
	if rl.IsKeyDown(rl.KeyLeftControl) && rl.IsKeyPressed(rl.KeyS) {
		b.Save()
	} else if b.Mode == Insert {
		if rl.IsKeyPressed(rl.KeyEscape) {
			b.Mode = Normal
		}
		if rl.IsKeyPressed(rl.KeyBackspace) {
			line := b.Lines[b.BufferCursor.Y]
			if b.BufferCursor.X > 0 {
				b.Lines[b.BufferCursor.Y] = line[:b.BufferCursor.X-1] + line[b.BufferCursor.X:]
				b.BufferCursor.X--
			} else if b.BufferCursor.Y > 0 {
				newBufferSlice := make([]string, 0)
				newBufferSlice = append(newBufferSlice, b.Lines[:b.BufferCursor.Y]...)
				newBufferSlice[b.BufferCursor.Y-1] += line
				newBufferSlice = append(newBufferSlice, b.Lines[b.BufferCursor.Y+1:]...)
				b.BufferCursor.Y--
				b.BufferCursor.X = utf8.RuneCountInString(b.Lines[b.BufferCursor.Y])
				b.Lines = newBufferSlice
			}
		}
		if rl.IsKeyPressed(rl.KeyEnter) {
			line := b.Lines[b.BufferCursor.Y]
			b.Lines[b.BufferCursor.Y] = line[:b.BufferCursor.X] + "\n" + line[b.BufferCursor.X:]
			newBufferSlice := make([]string, 0)
			newBufferSlice = append(newBufferSlice, b.Lines[:b.BufferCursor.Y]...)
			newBufferSlice = append(newBufferSlice, strings.Split(b.Lines[b.BufferCursor.Y], "\n")...)
			newBufferSlice = append(newBufferSlice, b.Lines[b.BufferCursor.Y+1:]...)
			b.Lines = newBufferSlice
			b.BufferCursor.Y++
			b.BufferCursor.X = 0
		}
		if rl.IsKeyPressed(rl.KeyDown) && b.BufferCursor.Y < len(b.Lines)-1 {
			b.BufferCursor.Y++
			wrapCursor(utf8.RuneCountInString(b.Lines[b.BufferCursor.Y-1]), b)
		}
		if rl.IsKeyPressed(rl.KeyUp) && b.BufferCursor.Y > 0 {
			b.BufferCursor.Y--
			wrapCursor(utf8.RuneCountInString(b.Lines[b.BufferCursor.Y+1]), b)
		}
		if rl.IsKeyPressed(rl.KeyRight) && b.BufferCursor.X < utf8.RuneCountInString(b.Lines[b.BufferCursor.Y]) {
			b.BufferCursor.X++
		}
		if rl.IsKeyPressed(rl.KeyLeft) && b.BufferCursor.X > 0 {
			b.BufferCursor.X--
		}
		key := rl.GetCharPressed()
		if key >= 32 && key <= 126 {
			line := b.Lines[b.BufferCursor.Y]
			b.Lines[b.BufferCursor.Y] = line[:b.BufferCursor.X] + string(rune(key)) + line[b.BufferCursor.X:]
			b.BufferCursor.X++
		}
	} else if b.Mode == Normal {
		if rl.IsKeyPressed(rl.KeyJ) && b.BufferCursor.Y < len(b.Lines)-1 {
			b.BufferCursor.Y++
			wrapCursor(utf8.RuneCountInString(b.Lines[b.BufferCursor.Y-1]), b)
		}
		if rl.IsKeyPressed(rl.KeyK) && b.BufferCursor.Y > 0 {
			b.BufferCursor.Y--
			wrapCursor(utf8.RuneCountInString(b.Lines[b.BufferCursor.Y+1]), b)
		}
		if rl.IsKeyPressed(rl.KeyL) && b.BufferCursor.X < utf8.RuneCountInString(b.Lines[b.BufferCursor.Y]) {
			b.BufferCursor.X++
		}
		if rl.IsKeyPressed(rl.KeyH) && b.BufferCursor.X > 0 {
			b.BufferCursor.X--
		}
		if rl.IsKeyPressed(rl.KeyI) {
			b.Mode = Insert
		}
	}
}

func (b *Buffer) Save() {
	file, err := os.Create(b.Path)
	if err != nil {
		LogError(err)
		return
	}
	content := strings.Join(b.Lines, "\n")
	_, err = file.WriteString(content)
	if err != nil {
		LogError(err)
	}
}
