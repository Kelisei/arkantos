package main

import (
	"arkantos/internal"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var theme internal.Theme
var config internal.Configuration
var openedBuffers []internal.Buffer
var currentBuffer int

// The start function, loads the configuration file, get's the theme selected
// setups the window, creates and initializes a buffer slice
func start(debug bool) error {
	var err error
	config, err = internal.LoadConfig("arkantos.conf")
	if err != nil {
		return err
	}

	theme, err = internal.ThemeParse(config.ThemeName)
	rl.InitWindow(800, 800, "Arkantos")
	rl.SetWindowState(rl.FlagWindowMaximized | rl.FlagWindowResizable)
	rl.MaximizeWindow()
	openedBuffers = make([]internal.Buffer, 0)
	currentBuffer = 0

	if debug {
		var debugBuffer internal.Buffer
		debugBuffer, err = internal.NewBuffer("testbuffer")
		if err != nil {
			return err
		}
		openedBuffers = append(openedBuffers, debugBuffer)
	}

	return err
}
func input() {

}
func update() {
}
func render() error {
	rl.BeginDrawing()
	rl.DrawFPS(1800, 10)
	rl.ClearBackground(theme.BgColor)
	err := openedBuffers[currentBuffer].RenderBuffer(config.FontSize, theme.FontColor)
	rl.EndDrawing()
	return err
}

func main() {
	err := start(true)
	defer rl.CloseWindow()
	if err != nil {
		fmt.Println(err)
		return
	}
	rl.SetTargetFPS(120)
	for !rl.WindowShouldClose() && err == nil {
		err = render()
	}
}
