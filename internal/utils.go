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
	ThemeName string
	FontSize  int
	MainFont  rl.Font
}

type jsonConfig struct {
	ThemeName string
	FontSize  int
	MainFont  string
}

type Theme struct {
	Name        string
	Version     string
	Description string
	BgColor     rl.Color
	FontColor   rl.Color
	Highlight   rl.Color
}

type jsonTheme struct {
	Name        string
	Version     string
	Description string
	BgColor     string
	FontColor   string
	Highlight   string
}

func hexToRaylib(hexa string) rl.Color {
	rgba, err := hex.DecodeString(hexa[1:])
	if err != nil {
		fmt.Println(err)
		return rl.White
	} else {
		return rl.NewColor(rgba[0], rgba[1], rgba[2], 255)
	}
}

func NewTheme(theme jsonTheme) Theme {
	t := Theme{}
	t.Name = theme.Name
	t.Version = theme.Version
	t.Description = theme.Description
	t.BgColor = hexToRaylib(theme.BgColor)
	t.FontColor = hexToRaylib(theme.FontColor)
	t.Highlight = hexToRaylib(theme.Highlight)
	return t
}

func NewConfig(c jsonConfig) Configuration {
	font := rl.LoadFont("static/fonts/" + c.MainFont)
	rl.SetTextureFilter(font.Texture, rl.FilterTrilinear)
	return Configuration{ThemeName: c.ThemeName, FontSize: c.FontSize, MainFont: font}
}

// Given the path string, it json decodes it and then returns a configuration struct
// with all the data.
func LoadConfig(path string) (Configuration, error) {

	file, err := os.Open(path)
	if err != nil {
		return Configuration{}, err
	}
	defer file.Close()
	config := jsonConfig{}
	err = json.NewDecoder(file).Decode(&config)
	return NewConfig(config), err
}

// Given the theme name, it json decodes it and then parses the values so
// it can be used by raylib.
func ThemeParse(theme_name string) (Theme, error) {
	file, _ := os.Open(theme_name + ".json")
	defer file.Close()
	theme := jsonTheme{}
	err := json.NewDecoder(file).Decode(&theme)
	if err != nil {
		return Theme{}, errors.New("error loading theme file")
	}
	return NewTheme(theme), nil
}
