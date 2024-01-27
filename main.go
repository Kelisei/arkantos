package main

import (
	"arkantos/internal"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var theme internal.Theme
var config internal.Configuration

func start() error {
	var err error
	config, err = internal.LoadConfig("arkantos.conf")
	if err != nil {
		return err
	}
	theme, err = internal.ThemeParse(config.ThemeName)
	rl.InitWindow(800, 800, "Arkantos")
	rl.SetWindowState(rl.FlagWindowMaximized | rl.FlagWindowResizable)
	rl.MaximizeWindow()
	return err
}
func input() {

}
func update() {

}
func render() {
	rl.BeginDrawing()
	rl.DrawFPS(10, 10)
	rl.ClearBackground(theme.BgColor)
	rl.EndDrawing()
}

func main() {
	err := start()
	defer rl.CloseWindow()
	if err != nil {
		fmt.Println(err)
		return
	}
	rl.SetTargetFPS(120)
	for !rl.WindowShouldClose() {
		render()
	}
}
