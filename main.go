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
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.InitWindow(1920, 1080, "Arkantos")
	var err error
	config, err = internal.LoadConfig("arkantos.conf")
	if err != nil {
		return err
	}
	theme, err = internal.ThemeParse(config.ThemeName)
	rl.SetWindowState(rl.FlagWindowMaximized | rl.FlagWindowResizable)
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.MaximizeWindow()
	rl.SetExitKey(rl.KeyF11)
	openedBuffers = make([]internal.Buffer, 0)
	currentBuffer = 0

	if debug {
		var debugBuffer internal.Buffer
		debugBuffer, err = internal.NewBuffer("testbuffer")
		if err != nil {
			return err
		}
		fmt.Printf("succesfully opened the testbuffer \n")
		openedBuffers = append(openedBuffers, debugBuffer)
	}

	return err
}
func input() {
	openedBuffers[currentBuffer].ListenInput()
}
func update() {
}
func render() error {
	rl.BeginDrawing()
	rl.DrawFPS(1800, 10)
	rl.ClearBackground(theme.BgColor)
	err := openedBuffers[currentBuffer].RenderBuffer(config.FontSize, theme.FontColor, theme.Highlight, config.MainFont)
	rl.EndDrawing()
	return err
}

func main() {
	err := start(true)
	if err != nil {
		internal.LogError(err)
		return
	}
	fmt.Println(openedBuffers[currentBuffer].BufferCursor)
	defer rl.CloseWindow()
	rl.SetTargetFPS(120)
	for !rl.WindowShouldClose() && err == nil {
		input()
		err = render()
		if err != nil {
			internal.LogError(err)
		}
	}
}
