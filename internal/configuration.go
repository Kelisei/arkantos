package internal

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	rl "github.com/gen2brain/raylib-go/raylib"
)

type configurationToml struct {
	Version             string
	FontSize            float32
	TargetFPS           int32
	IndentSize          int
	FontRegular         string
	FontBold            string
	FontItalic          string
	PrimaryColor        string
	SecondaryColor      string
	TertiaryColor       string
	AccentColor         string
	AccentColor2        string
	DefaultFontColor    string
	HighlightFontColor  string
	SecondaryFontColor  string
	KeywordsFontColor   string
	FunctionsFontColor  string
	VariablesFontColor  string
	LiteralsFontColor   string
	RelativeLineNumbers bool
}

type Configuration struct {
	Version             string
	FontSize            float32
	TargetFPS           int32
	IndentSize          int
	FontRegular         rl.Font
	FontBold            rl.Font
	FontItalic          rl.Font
	PrimaryColor        rl.Color
	SecondaryColor      rl.Color
	TertiaryColor       rl.Color
	AccentColor         rl.Color
	AccentColor2        rl.Color
	DefaultFontColor    rl.Color
	HighlightFontColor  rl.Color
	SecondaryFontColor  rl.Color
	KeywordsFontColor   rl.Color
	FunctionsFontColor  rl.Color
	VariablesFontColor  rl.Color
	LiteralsFontColor   rl.Color
	RelativeLineNumbers bool
}

func NewConfiguration(path string) (Configuration, error) {
	configPath, err := filepath.Abs(path)
	if err != nil {
		return Configuration{}, err
	}
	var file *os.File
	file, err = os.Open(configPath)
	if err != nil {
		return Configuration{}, err
	}
	var config configurationToml
	_, err = toml.NewDecoder(file).Decode(&config)
	if err != nil {
		return Configuration{}, err
	}
	return parseConfiguration(config), nil
}

// The most important part of this is that we make raylib create a font atlas where each glyph is 80px in height.
func parseConfiguration(c configurationToml) Configuration {
	fontReg := rl.LoadFontEx(c.FontRegular, 80, nil)
	fontItalic := rl.LoadFontEx(c.FontItalic, 80, nil)
	fontBold := rl.LoadFontEx(c.FontBold, 80, nil)
	rl.SetTextureFilter(fontBold.Texture, rl.FilterTrilinear)
	rl.SetTextureFilter(fontItalic.Texture, rl.FilterBilinear)
	rl.SetTextureFilter(fontReg.Texture, rl.FilterTrilinear)
	return Configuration{
		Version: c.Version, FontSize: c.FontSize,
		TargetFPS:           c.TargetFPS,
		IndentSize:          c.IndentSize,
		FontRegular:         fontReg,
		FontBold:            fontBold,
		FontItalic:          fontItalic,
		PrimaryColor:        hexToRaylib(c.PrimaryColor),
		SecondaryColor:      hexToRaylib(c.SecondaryColor),
		TertiaryColor:       hexToRaylib(c.TertiaryColor),
		AccentColor:         hexToRaylib(c.AccentColor),
		AccentColor2:        hexToRaylib(c.AccentColor2),
		DefaultFontColor:    hexToRaylib(c.DefaultFontColor),
		HighlightFontColor:  hexToRaylib(c.HighlightFontColor),
		SecondaryFontColor:  hexToRaylib(c.SecondaryFontColor),
		KeywordsFontColor:   hexToRaylib(c.KeywordsFontColor),
		FunctionsFontColor:  hexToRaylib(c.FunctionsFontColor),
		VariablesFontColor:  hexToRaylib(c.VariablesFontColor),
		LiteralsFontColor:   hexToRaylib(c.LiteralsFontColor),
		RelativeLineNumbers: c.RelativeLineNumbers,
	}
}
