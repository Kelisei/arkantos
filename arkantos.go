package main

import (
	"arkantos/internal"
	"errors"
	"fmt"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var theme internal.Theme
var config internal.Configuration
var openedBuffers []internal.Buffer
var currentBuffer int
var closeWindow bool

// The start function, loads the configuration file, get's the theme selected
// setups the window, creates and initializes a buffer slice
func start() error {
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.InitWindow(1920, 1080, "Arkantos")
	var err error
	config, err = internal.LoadConfig("arkantos.conf")
	if err != nil {
		return err
	}
	theme, err = internal.ThemeParse(config.ThemeName)
	if err != nil {
		return err
	}
	rl.SetWindowState(rl.FlagWindowMaximized | rl.FlagWindowResizable)
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	rl.MaximizeWindow()
	rl.SetExitKey(rl.KeyF11)
	openedBuffers = make([]internal.Buffer, 0)
	currentBuffer = 0
	args := os.Args
	fmt.Println(args)
	if len(args) < 2 || len(args) > 2 {
		return errors.New("wrong number of parameters, did you mean to pass a file path?")
	}
	var buffer internal.Buffer
	buffer, err = internal.NewBuffer(args[1])
	if err != nil {
		return err
	}
	openedBuffers = append(openedBuffers, buffer)
	return err
}
func input() {
	openedBuffers[currentBuffer].ListenInput(&closeWindow, config.MainFont, config.FontSize)
}
func render() {
	rl.BeginDrawing()
	rl.DrawFPS(1800, 10)
	rl.ClearBackground(theme.BgColor)
	openedBuffers[currentBuffer].RenderBuffer(config.FontSize, theme.FontColor, theme.Highlight, theme.BottomBarColor, theme.BottomBarFontColor, config.MainFont)
	rl.EndDrawing()
}

func main() {
	err := start()
	if err != nil {
		internal.LogError(err)
		return
	}
	defer rl.CloseWindow()
	rl.SetTargetFPS(120)
	for !rl.WindowShouldClose() && !closeWindow {
		input()
		render()
	}
}
