package internal

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	sitter "github.com/smacker/go-tree-sitter"

	"github.com/smacker/go-tree-sitter/golang"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const (
	INSERT = iota
	NORMAL
)

type Change struct {
	StartingLine int
	Lines        []string
}

type Buffer struct {
	Path           string
	Content        string
	BCursor        Cursor
	Mode           int
	History        []Change
	CurrentCommand string
	CurrentChange  int
	YScroll        int32
	XScroll        int32
	X              int32
	Y              int32
	Width          int32
	Height         int32
	Exists         bool
	Lang           *sitter.Language
}

func NewBufferFromPath(path string, langs map[string]*sitter.Language) (Buffer, error) {
	bufferPath, errPath := filepath.Abs(path)
	if errPath != nil {
		return Buffer{}, errors.New("Failed to find path to buffer: " + errPath.Error())
	}
	content, errRead := os.ReadFile(bufferPath)
	if errRead != nil {
		return Buffer{}, errors.New("Failed to read from buffer: " + errRead.Error())
	}
	cursor := Cursor{X: 0, Y: 0, Length: 0}
	split := strings.Split(bufferPath, ".")
	suffix := split[len(split)-1]

	return Buffer{
		Path:           bufferPath,
		Content:        string(content),
		BCursor:        cursor,
		Mode:           INSERT,
		History:        make([]Change, 0),
		CurrentCommand: "",
		CurrentChange:  0,
		YScroll:        0,
		XScroll:        0,
		Exists:         true,
		Lang:           langs[suffix],
	}, nil
}

func NewEmptyBuffer() (Buffer, error) {
	path, err := filepath.Abs("blank-name")
	if err != nil {
		return Buffer{}, err
	}
	return Buffer{
		Path:           path,
		Content:        "",
		BCursor:        Cursor{X: 0, Y: 0, Length: 0},
		Mode:           INSERT,
		CurrentCommand: "",
		CurrentChange:  0,
		YScroll:        0,
		XScroll:        0,
		Exists:         false,
		Lang:           nil,
	}, nil
}

func (b *Buffer) Update() {
	file, err := os.Open(b.Path)
	if err != nil {
		b.Exists = false
	} else {
		file.Close()
	}
}

func (b *Buffer) Render(c Configuration, langs map[string]*sitter.Language) {
	rl.BeginScissorMode(b.X, b.Y, b.Width, b.Height)
	switch b.Lang {
	default:
		lines := strings.Split(b.Content, "\n")
		b.BCursor.draw(
			int(c.FontSize*2),
			b.YScroll,
			b.XScroll,
			c.FontSize,
			lines[b.BCursor.Y],
			c.FontRegular,
			c.DefaultFontColor,
		)
		for i, line := range lines {
			lineNumber := i
			var lineNumberColor rl.Color
			if i == b.BCursor.Y {
				lineNumberColor = c.DefaultFontColor
			} else {
				lineNumberColor = c.SecondaryFontColor
			}
			if c.RelativeLineNumbers {
				if lineNumber != b.BCursor.Y {
					lineNumber = absInt(i-b.BCursor.Y) - 1
				}
			}
			lineYPos := float32(b.Y) + (c.FontSize * float32(i-int(b.YScroll)))
			rl.DrawTextEx(
				c.FontRegular,
				strconv.Itoa(lineNumber+1),
				rl.NewVector2(float32(b.XScroll), lineYPos),
				c.FontSize,
				0,
				lineNumberColor,
			)
			rl.DrawTextEx(
				c.FontRegular,
				line,
				rl.NewVector2(c.FontSize*2, lineYPos),
				c.FontSize,
				0,
				c.DefaultFontColor,
			)
		}
	}
	b.DrawBottomBar(c)
	rl.EndScissorMode()
}

func (b *Buffer) DrawBottomBar(c Configuration) {
	rl.DrawRectangle(
		0,
		b.Height-int32(c.FontSize*2),
		b.Width,
		int32(c.FontSize*2),
		c.TertiaryColor,
	)
	cmd := ":" + b.CurrentCommand
	mode := ""
	var color *rl.Color
	switch b.Mode {
	case INSERT:
		mode = "INSERT"
		color = &c.AccentColor
	case NORMAL:
		mode = "NORMAL"
		color = &c.AccentColor2
	}
	box := rl.MeasureTextEx(c.FontBold, mode, c.FontSize, 0)
	rl.DrawRectangle(0, b.Height-int32(c.FontSize*2), int32(box.X), int32(box.Y), *color)
	rl.DrawTextEx(
		c.FontBold,
		mode,
		rl.NewVector2(0, float32(b.Height-int32(c.FontSize*2))),
		c.FontSize,
		0,
		c.SecondaryFontColor,
	)
	rl.DrawTextEx(
		c.FontRegular,
		cmd,
		rl.NewVector2(box.X, float32(b.Height-int32(c.FontSize*2))),
		c.FontSize,
		0,
		c.DefaultFontColor,
	)
	rl.DrawTextEx(
		c.FontItalic,
		b.Path,
		rl.NewVector2(0, float32(b.Height-int32(c.FontSize))),
		c.FontSize,
		0,
		c.DefaultFontColor,
	)
}

func (b *Buffer) Save() {
	file, err := os.Create(b.Path)
	if err != nil {
		b.CurrentCommand = "Failed to overwrite file"
		LogError(err)
		return
	}
	_, err = file.WriteString(b.Content)
	if err != nil {
		b.CurrentCommand = "Failed to write into file"
		LogError(err)
	}
	b.CurrentCommand = "File saved"
}

func getLinesAndUpdate(content string) ([]string, int64) {
	return strings.Split(content, "\n"), time.Now().UnixMilli()
}

func insertLines(lines []string, y int, linesToInsert ...string) []string {
	content := make([]string, 0)
	content = append(content, lines[:y]...)
	for _, line := range linesToInsert {
		content = append(content, line)
	}
	content = append(content, lines[y+1:]...)
	return content
}

func (buffer *Buffer) ListenInput(config Configuration, state *State, closeWindow *bool) {
	if rl.IsKeyDown(rl.KeyLeftControl) && rl.IsKeyDown(rl.KeyS) {
		buffer.Save()
	}
	if buffer.Mode == INSERT {
		if rl.IsKeyPressed(rl.KeyDown) || (rl.IsKeyDown(rl.KeyDown) && IsUpdateTick(*state)) {
			getLinesAndMoveY(1, state, buffer, config)
		}
		if rl.IsKeyPressed(rl.KeyUp) || (rl.IsKeyDown(rl.KeyUp) && IsUpdateTick(*state)) {
			getLinesAndMoveY(-1, state, buffer, config)
		}
		if rl.IsKeyPressed(rl.KeyRight) || (rl.IsKeyDown(rl.KeyRight) && IsUpdateTick(*state)) {
			getLinesAndMoveX(1, state, buffer)
		}
		if rl.IsKeyPressed(rl.KeyLeft) || (rl.IsKeyDown(rl.KeyRight) && IsUpdateTick(*state)) {
			getLinesAndMoveX(-1, state, buffer)
		}
		if rl.IsKeyPressed(rl.KeyEnter) || (rl.IsKeyDown(rl.KeyEnter) && IsUpdateTick(*state)) {
			var lines []string
			lines, state.lastUpdateTime = getLinesAndUpdate(buffer.Content)
			line1, line2 := lines[buffer.BCursor.Y][:buffer.BCursor.X], lines[buffer.BCursor.Y][buffer.BCursor.X:]
			buffer.Content = strings.Join(insertLines(lines, buffer.BCursor.Y, line1, line2), "\n")
			buffer.BCursor.moveY(1, len(lines), lines, config.FontSize, &buffer.YScroll, buffer.Height)
			buffer.BCursor.X = 0
		}
		if rl.IsKeyPressed(rl.KeyBackspace) || (rl.IsKeyDown(rl.KeyBackspace) && IsUpdateTick(*state)) {
			var lines []string
			lines, state.lastUpdateTime = getLinesAndUpdate(buffer.Content)
			if buffer.BCursor.X > 0 {
				lines[buffer.BCursor.Y] = lines[buffer.BCursor.Y][:buffer.BCursor.X-1] + lines[buffer.BCursor.Y][buffer.BCursor.X:]
				buffer.BCursor.X--
			} else if buffer.BCursor.X == 0 && buffer.BCursor.Y > 0 {
				previousLineLen := utf8.RuneCountInString(lines[buffer.BCursor.Y-1])
				lines[buffer.BCursor.Y-1] += lines[buffer.BCursor.Y]
				lines = append(lines[:buffer.BCursor.Y], lines[buffer.BCursor.Y+1:]...)
				buffer.BCursor.moveY(-1, len(lines), lines, config.FontSize, &buffer.YScroll, buffer.Height)
				buffer.BCursor.X = previousLineLen
			}
			buffer.Content = strings.Join(lines, "\n")
		}
		if rl.IsKeyPressed(rl.KeyEscape) {
			buffer.Mode = NORMAL
		}
		key := rl.GetCharPressed()
		if key >= 32 && key <= 126 {
			lines := strings.Split(buffer.Content, "\n")
			line := lines[buffer.BCursor.Y]
			lines[buffer.BCursor.Y] = line[:buffer.BCursor.X] + string(rune(key)) + line[buffer.BCursor.X:]
			buffer.BCursor.moveX(1, utf8.RuneCountInString(lines[buffer.BCursor.Y]))
			buffer.Content = strings.Join(lines, "\n")
		}
	} else if buffer.Mode == NORMAL {
		if rl.IsKeyPressed(rl.KeyTab) {
			state.cycleBuffer()
		}
		if rl.IsKeyPressed(rl.KeyI) {
			buffer.Mode = INSERT
			buffer.CurrentCommand = "i"
		}
		if rl.IsKeyPressed(rl.KeyJ) || (rl.IsKeyDown(rl.KeyJ) && IsUpdateTick(*state)) {
			buffer.CurrentCommand = "j"
			getLinesAndMoveY(1, state, buffer, config)
		}
		if rl.IsKeyPressed(rl.KeyK) || (rl.IsKeyDown(rl.KeyK) && IsUpdateTick(*state)) {
			buffer.CurrentCommand = "k"
			getLinesAndMoveY(-1, state, buffer, config)
		}
		if rl.IsKeyPressed(rl.KeyL) || (rl.IsKeyDown(rl.KeyL) && IsUpdateTick(*state)) {
			buffer.CurrentCommand = "l"
			getLinesAndMoveX(1, state, buffer)
		}
		if rl.IsKeyPressed(rl.KeyH) || (rl.IsKeyDown(rl.KeyH) && IsUpdateTick(*state)) {
			buffer.CurrentCommand = "h"
			getLinesAndMoveX(-1, state, buffer)
		}
		if rl.IsKeyPressed(rl.KeyW) {
			buffer.CurrentCommand = "w"
		}
		if rl.IsKeyPressed(rl.KeyQ) {
			buffer.CurrentCommand += "q"
		}
		if rl.IsKeyPressed(rl.KeyEnter) {
			if buffer.CurrentCommand == "w" {
				buffer.Save()
			} else if buffer.CurrentCommand == "q" {
				*closeWindow = true
			} else if buffer.CurrentCommand == "wq" {
				buffer.Save()
				*closeWindow = true
			}
		}
	}
}

func getLinesAndMoveX(amount int, state *State, buffer *Buffer) {
	var lines []string
	lines, state.lastUpdateTime = getLinesAndUpdate(buffer.Content)
	buffer.BCursor.moveX(amount, utf8.RuneCountInString(lines[buffer.BCursor.Y]))
}

func getLinesAndMoveY(amount int, s *State, b *Buffer, c Configuration) {
	var lines []string
	lines, s.lastUpdateTime = getLinesAndUpdate(b.Content)
	b.BCursor.moveY(amount, len(lines), lines, c.FontSize, &b.YScroll, b.Height)
}

func GetLanguages() map[string]*sitter.Language {
	langs := make(map[string]*sitter.Language)
	langs["go"] = golang.GetLanguage()
	return langs
}
