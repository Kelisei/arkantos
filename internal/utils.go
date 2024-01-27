package internal

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Configuration struct {
	ThemeName  string
	FontSize   int
	FontFamily string
}

type Theme struct {
	Name        string
	Version     string
	Description string
	BgColor     rl.Color
	FontColor   rl.Color
}

type json_theme struct {
	Name        string
	Version     string
	Description string
	BgColor     string
	FontColor   string
}

func NewTheme(theme json_theme) Theme {
	t := Theme{}
	t.Name = theme.Name
	t.Version = theme.Version
	t.Description = theme.Description
	rgba, err := hex.DecodeString(theme.BgColor[1:])
	if err != nil {
		fmt.Println(err)
		t.BgColor = rl.White
	} else {
		t.BgColor = rl.NewColor(rgba[0], rgba[1], rgba[2], 255)
	}
	rgba, err = hex.DecodeString(theme.FontColor[1:])
	if err != nil {
		fmt.Println(err)
		t.FontColor = rl.Blue
	} else {
		t.FontColor = rl.NewColor(rgba[0], rgba[1], rgba[2], 255)
	}
	return t
}

// Given the path string, it json decodes it and then returns a configuration struct
// with all the data.
func LoadConfig(path string) (Configuration, error) {
	config := Configuration{}
	file, err := os.Open(path)
	if err != nil {
		return config, err
	}
	defer file.Close()
	err = json.NewDecoder(file).Decode(&config)
	return config, err
}

// Given the theme name, it json decodes it and then parses the values so
// it can be used by raylib.
func ThemeParse(theme_name string) (Theme, error) {
	file, _ := os.Open(theme_name + ".json")
	defer file.Close()
	theme := json_theme{}
	err := json.NewDecoder(file).Decode(&theme)
	if err != nil {
		return Theme{}, errors.New("error loading theme file")
	}
	return NewTheme(theme), nil
}
