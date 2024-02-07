package main

import (
	"arkantos/internal"
	"fmt"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
	window := internal.NewWindow(1080, 1920, rl.KeyF11, true, true, "ARKANTOS")
	internal.StartWindow(&window)
	defer rl.CloseWindow()

	config, err := internal.NewConfiguration("config.toml")
	if err != nil {
		internal.LogError(err)
		return
	}
	rl.SetTargetFPS(config.TargetFPS)
	state := internal.State{}
	langs := internal.GetLanguages()
	args := os.Args
	if len(args) > 1 {
		for _, arg := range args[1:] {
			buff, err := internal.NewBufferFromPath(arg, langs)
			if err != nil {
				internal.LogError(err)
			} else {
				state.Buffers = append(state.Buffers, buff)
			}
		}
	} else {
		buff, err := internal.NewEmptyBuffer()
		if err != nil {
			internal.LogError(err)
		} else {
			state.Buffers = append(state.Buffers, buff)
		}
	}
	if len(state.Buffers) == 0 {
		internal.LogString("no buffers were able to be open")
		return
	}
	fmt.Println(state.Buffers[state.CurrentBuffer].Path)
	for !window.CloseWindow {
		window.CloseWindow = rl.WindowShouldClose()
		rl.BeginDrawing()
		rl.ClearBackground(config.PrimaryColor)
		currentBuffer := &state.Buffers[state.CurrentBuffer]
		currentBuffer.Width, currentBuffer.Height = internal.GetWindowSize(&window)
		currentBuffer.Render(config, langs)
		currentBuffer.ListenInput(config, &state, &window.CloseWindow)
		if len(state.Buffers) == 0 {
			window.CloseWindow = true
		}
		rl.EndDrawing()
	}
}
