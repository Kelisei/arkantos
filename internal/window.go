package internal

import rl "github.com/gen2brain/raylib-go/raylib"

// The window struct acts as window as a interface between the raylib window and the application, please take into account:
// The height and width of this window is only for the starting size if the window is resizable.
type Window struct {
	Width          int32
	Height         int32
	ExitKey        int32
	Resizable      bool
	StartMaximized bool
	CloseWindow    bool
	Name           string
	Buffer         []Buffer
}

func NewWindow(
	width, height, exitkey int32,
	resizable, startMaximized bool,
	name string,
) Window {
	return Window{
		Width:          width,
		Height:         height,
		ExitKey:        exitkey,
		Resizable:      resizable,
		StartMaximized: startMaximized,
		CloseWindow:    false,
		Name:           name,
	}
}

func StartWindow(w *Window) {
	rl.InitWindow(w.Width, w.Height, w.Name)
	rl.SetConfigFlags(rl.FlagMsaa4xHint)
	if w.Resizable {
		rl.SetWindowState(rl.FlagWindowResizable)
	}
	if w.StartMaximized {
		rl.MaximizeWindow()
	}
	rl.SetExitKey(w.ExitKey)
}

func GetWindowSize(w *Window) (int32, int32) {
	if !w.Resizable {
		return w.Width, w.Height
	}
	return int32(rl.GetScreenWidth()), int32(rl.GetScreenHeight())
}
func GetWindowSizeFloat(w *Window) (float32, float32) {
	if !w.Resizable {
		return float32(w.Width), float32(w.Width)
	}
	return float32(rl.GetScreenWidth()), float32(rl.GetScreenHeight())
}
