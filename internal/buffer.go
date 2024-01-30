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

type Buffer struct {
	Path           string
	Mode           int
	Lines          []string
	BCursor        Cursor
	UndoList       []Change
	Padding        int
	CurrentCommand string
}

// Given a path, it creates a new Buffer with default values and the file's content.
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
	return Buffer{Path: path, Mode: Insert, Lines: lines, BCursor: cursor, UndoList: undos, Padding: 10, CurrentCommand: ""}, nil
}

// Draws in screen the lines in order, and the bottom info bar.
func (b *Buffer) RenderBuffer(fontSize int, fontColor, highlight, bottomBarColor, bottomBarFontColor rl.Color, font rl.Font) {
	drawLines(b, fontSize, font, fontColor, highlight)
	drawBottomBar(fontSize, b, font, bottomBarFontColor, bottomBarColor)
}

func drawLines(b *Buffer, fontSize int, font rl.Font, fontColor, highlight rl.Color) {
	// YOffset is how much the lines should be moved up or down in order to see them.
	b.BCursor.DrawCursor(fontSize, b.Padding, highlight, *b, font)
	for i, line := range b.Lines {
		var lineNumber string
		lineNumberColor := rl.White
		if i != b.BCursor.Y {
			lineNumber = strconv.Itoa(absInt(i - b.BCursor.Y))
		} else {
			lineNumber = strconv.Itoa(i + 1)
			lineNumberColor = highlight
		}
		position := rl.Vector2{X: float32(b.Padding + 1), Y: float32((i + b.BCursor.YOffset) * fontSize)}
		rl.DrawTextEx(font, lineNumber, position, float32(fontSize), 0, lineNumberColor)
		position.X = float32(b.Padding + fontSize*2)
		rl.DrawTextEx(font, line, position, float32(fontSize), 0, fontColor)
	}
}

// Given certain parameters and the size of the screen, draws the current buffer,
// the cursor position and the current mode.
func drawBottomBar(fontSize int, b *Buffer, font rl.Font, fontColor, bg rl.Color) {
	height := rl.GetScreenHeight()
	width := rl.GetScreenWidth()
	bottomPos := float32(height - fontSize)

	rl.DrawRectangle(0, int32(bottomPos-float32(fontSize)), int32(width), int32(fontSize*2), bg)

	statusPos := rl.Vector2{X: float32(b.Padding), Y: bottomPos - float32(fontSize)}
	statusString := ":" + b.CurrentCommand
	rl.DrawTextEx(font, statusString, statusPos, float32(fontSize), 0, fontColor)

	pathPos := rl.Vector2{X: float32(b.Padding), Y: bottomPos}
	pathStr := "Buffer :" + b.Path
	rl.DrawTextEx(font, pathStr, pathPos, float32(fontSize), 0, fontColor)

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
	rl.DrawTextEx(font, "--"+currentMode+"--", modePos, float32(fontSize), 0, fontColor)

	cursorPosStr := strconv.Itoa(b.BCursor.X+1) + "," + strconv.Itoa(b.BCursor.Y+1)
	cursorPos := rl.Vector2{X: float32(rl.MeasureText(pathStr, int32(fontSize))), Y: bottomPos}
	rl.DrawTextEx(font, cursorPosStr, cursorPos, float32(fontSize), 0, fontColor)
}

// Listens for input and acts on it, allows for complex commands.
func (b *Buffer) ListenInput(closeWindow *bool, font rl.Font, fontSize int) {
	if rl.IsKeyDown(rl.KeyLeftControl) && rl.IsKeyPressed(rl.KeyS) {
		b.Save()
	} else if b.Mode == Insert {
		if rl.IsKeyPressed(rl.KeyEscape) {
			b.Mode = Normal
			b.CurrentCommand = "ESC"
		}
		if rl.IsKeyPressed(rl.KeyBackspace) {
			line := b.Lines[b.BCursor.Y]
			if b.BCursor.X > 0 {
				b.Lines[b.BCursor.Y] = line[:b.BCursor.X-1] + line[b.BCursor.X:]
				b.BCursor.X--
			} else if b.BCursor.Y > 0 {
				newBufferSlice := make([]string, 0)
				newBufferSlice = append(newBufferSlice, b.Lines[:b.BCursor.Y]...)
				newBufferSlice[b.BCursor.Y-1] += line
				newBufferSlice = append(newBufferSlice, b.Lines[b.BCursor.Y+1:]...)
				b.BCursor.Y--
				b.BCursor.X = utf8.RuneCountInString(b.Lines[b.BCursor.Y])
				b.Lines = newBufferSlice
			}
		}
		if rl.IsKeyPressed(rl.KeyEnter) {
			line := b.Lines[b.BCursor.Y]
			b.Lines[b.BCursor.Y] = line[:b.BCursor.X] + "\n" + line[b.BCursor.X:]
			newBufferSlice := make([]string, 0)
			newBufferSlice = append(newBufferSlice, b.Lines[:b.BCursor.Y]...)
			newBufferSlice = append(newBufferSlice, strings.Split(b.Lines[b.BCursor.Y], "\n")...)
			newBufferSlice = append(newBufferSlice, b.Lines[b.BCursor.Y+1:]...)
			b.Lines = newBufferSlice
			b.BCursor.Y++
			b.BCursor.X = 0
		}
		if rl.IsKeyPressed(rl.KeyDown) && b.BCursor.Y < len(b.Lines)-1 {
			b.BCursor.Y++
			b.BCursor.WrapCursor(utf8.RuneCountInString(b.Lines[b.BCursor.Y-1]), b)
			lineHeight := rl.MeasureTextEx(font, " ", float32(fontSize), 0).Y
			if lineHeight*float32(b.BCursor.Y-b.BCursor.YOffset) >= float32(rl.GetScreenHeight()-fontSize*2) {
				b.BCursor.ChangeYOffset(-1)
			}
		}
		if rl.IsKeyPressed(rl.KeyUp) && b.BCursor.Y > 0 {
			b.BCursor.Y--
			b.BCursor.WrapCursor(utf8.RuneCountInString(b.Lines[b.BCursor.Y+1]), b)
			if (b.BCursor.Y+b.BCursor.YOffset)*fontSize < 0 {
				b.BCursor.ChangeYOffset(1)
			}
		}
		if rl.IsKeyPressed(rl.KeyRight) && b.BCursor.X < utf8.RuneCountInString(b.Lines[b.BCursor.Y]) {
			b.BCursor.X++
		}
		if rl.IsKeyPressed(rl.KeyLeft) && b.BCursor.X > 0 {
			b.BCursor.X--
		}
		key := rl.GetCharPressed()
		if key >= 32 && key <= 126 {
			line := b.Lines[b.BCursor.Y]
			b.Lines[b.BCursor.Y] = line[:b.BCursor.X] + string(rune(key)) + line[b.BCursor.X:]
			b.BCursor.X++
		}
	} else if b.Mode == Normal {
		if rl.IsKeyPressed(rl.KeyW) {
			b.CurrentCommand = "w"
		}
		if rl.IsKeyPressed(rl.KeyQ) {
			if b.CurrentCommand == "w" {
				b.CurrentCommand += "q"
			}
			b.CurrentCommand = "q"
		}
		if rl.IsKeyPressed(rl.KeyEnter) {
			switch b.CurrentCommand {
			case "w":
				b.Save()
			case "q":
				*closeWindow = true
			case "wq":
				b.Save()
				*closeWindow = true
			}
		}
		if rl.IsKeyPressed(rl.KeyJ) && b.BCursor.Y < len(b.Lines)-1 {
			b.BCursor.Y++
			b.BCursor.WrapCursor(utf8.RuneCountInString(b.Lines[b.BCursor.Y-1]), b)
			b.CurrentCommand = "j"
		}
		if rl.IsKeyPressed(rl.KeyK) && b.BCursor.Y > 0 {
			b.BCursor.Y--
			b.BCursor.WrapCursor(utf8.RuneCountInString(b.Lines[b.BCursor.Y+1]), b)
			b.CurrentCommand = "k"
		}
		if rl.IsKeyPressed(rl.KeyL) && b.BCursor.X < utf8.RuneCountInString(b.Lines[b.BCursor.Y]) {
			b.BCursor.X++
			b.CurrentCommand = "l"
		}
		if rl.IsKeyPressed(rl.KeyH) && b.BCursor.X > 0 {
			b.BCursor.X--
			b.CurrentCommand = "h"
		}
		if rl.IsKeyPressed(rl.KeyI) {
			b.Mode = Insert
			b.CurrentCommand = "i"
		}

	}
}

// Overwrite's the file with the content of the buffer.
func (b *Buffer) Save() {
	file, err := os.Create(b.Path)
	if err != nil {
		b.CurrentCommand = "Failed to overwrite file"
		LogError(err)
		return
	}
	content := strings.Join(b.Lines, "\n")
	_, err = file.WriteString(content)
	if err != nil {
		b.CurrentCommand = "Failed to write into file"
		LogError(err)
	}
	b.CurrentCommand = "Buffer saved"
}
