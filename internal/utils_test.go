package internal

import (
	"strconv"
	"testing"
)

func TestLoadConfig(t *testing.T) {

	t.Run("valid config", func(t *testing.T) {
		const validPath = "testdata/validconfig.json"

		conf, err := LoadConfig(validPath)

		if err != nil {
			t.Fatal(err)
		}
		if conf.ThemeName != "darktheme" {
			t.Fatal("Expected conf.ThemeName to be 'darktheme' but got " + conf.ThemeName)
		}
		if conf.FontSize != 14 {
			t.Fatal("Expected conf.FontSize to be '14' but got " + strconv.Itoa(conf.FontSize))
		}
	})

	t.Run("invalid path", func(t *testing.T) {
		const invalidPath = "testdata/does_not_exist.json"

		_, err := LoadConfig(invalidPath)

		if err == nil {
			t.Error("Expected error for invalid path but got none")
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		const invalidJSONPath = "testdata/invalidconfig.json"

		_, err := LoadConfig(invalidJSONPath)

		if err == nil {
			t.Error("Expected error for invalid JSON but got none")
		}
	})

}

func TestThemeParse(t *testing.T) {

	t.Run("valid theme", func(t *testing.T) {
		theme, err := ThemeParse("testdata/validtheme")
		if err != nil {
			t.Fatal(err)
		}
		if theme.Name != "darktheme.json" {
			t.Errorf("Expected name to be 'darktheme.json' but got %s", theme.Name)
		}
	})

	t.Run("invalid theme file", func(t *testing.T) {
		_, err := ThemeParse("testdata/invalidtheme")
		if err == nil {
			t.Error("Expected error but got none")
		}
	})

	t.Run("theme not found", func(t *testing.T) {
		_, err := ThemeParse("testdata/doesnotexist")
		if err == nil {
			t.Errorf("Expected not exists error but got %v", err)
		}
	})

}
