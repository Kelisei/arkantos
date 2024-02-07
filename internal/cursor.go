package internal

import (
	"fmt"
	"unicode/utf8"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Cursor struct {
	Y      int
	X      int
	Length int
}

func (c *Cursor) moveY(
	amount, bufferLen int,
	lines []string,
	fontSize float32,
	yScroll *int32,
	bufferHeight int32,
) {
	if c.Y+amount < bufferLen && c.Y+amount >= 0 {
		c.Y += amount
		if utf8.RuneCountInString(lines[c.Y]) < c.X {
			c.X = utf8.RuneCountInString(lines[c.Y])
			fmt.Println("new x positions", c.X)
		}
		if fontSize*float32(int32(c.Y)-*yScroll) >= float32(bufferHeight) {
			*yScroll += int32(amount)
		} else if fontSize*float32(int32(c.Y)-*yScroll) < 0 {
			*yScroll += int32(amount)
		}
	}
}

func (c *Cursor) moveX(
	amount, lineLen int,
) {
	c.X += amount
	if c.X < 0 {
		c.X = 0
	} else if c.X > lineLen {
		c.X = lineLen
	}
}

func (c *Cursor) draw(
	lineStart int,
	scrollY, scrollX int32,
	fontSize float32,
	currentLine string,
	font rl.Font,
	color rl.Color,
) {
	sliced := currentLine[:c.X]
	measurement := rl.MeasureTextEx(font, sliced, fontSize, 0)
	x := int32(lineStart+int(measurement.X)) + scrollX
	yStart := int32(fontSize) * (int32(c.Y) - scrollY)
	yEnd := yStart + int32(fontSize)

	rl.DrawLine(x, yStart, x, yEnd, color)
}
