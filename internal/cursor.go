package internal

import (
	"unicode/utf8"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Cursor struct {
	X       int
	Y       int
	XOffset int
	YOffset int
}

func (c *Cursor) DrawCursor(fontSize, padding int, color rl.Color, b Buffer, font rl.Font) {
	lineStart := padding + fontSize*2
	line := b.Lines[c.Y]

	xStart := int32(lineStart + int(rl.MeasureTextEx(font, line[:c.X], float32(fontSize), 0).X))
	yStart := int32(fontSize * (c.Y + c.YOffset))
	xEnd := xStart
	yEnd := yStart + int32(fontSize+c.YOffset)

	rl.DrawLine(xStart, yStart, xEnd, yEnd, color)
}

func (c *Cursor) WrapCursor(previousLineLen int, b *Buffer) {
	if c.X > utf8.RuneCountInString(b.Lines[c.Y]) || c.X == previousLineLen {
		c.X = utf8.RuneCountInString(b.Lines[c.Y])
	}
}

func (c *Cursor) ChangeYOffset(amount int) {
	if c.YOffset+amount <= 0 {
		c.YOffset += amount
	}
}
