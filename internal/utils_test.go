package internal

import (
	"strconv"
	"testing"
)

func TestLoadConfig(t *testing.T) {

	t.Run("valid config", func(t *testing.T) {
		// Arrange
		const validPath = "testdata/valid.json"

		// Act
		conf, err := LoadConfig(validPath)

		// Assert
		if err != nil {
			t.Fatal(err)
		}
		// assert conf is populated as expected
		if conf.ThemeName != "darktheme" {
			t.Fatal("Expected conf.ThemeName to be 'darktheme' but got " + conf.ThemeName)
		}
		if conf.FontSize != 14 {
			t.Fatal("Expected conf.FontSize to be '14' but got " + strconv.Itoa(conf.FontSize))
		}
		if conf.FontFamily != "Consolas, Menlo, Monaco, 'Courier New', monospace" {
			t.Fatal("Expected conf.FontFamily to be 'Consolas, Menlo, Monaco, 'Courier New', monospace' but got " + conf.FontFamily)
		}

	})

	t.Run("invalid path", func(t *testing.T) {
		// Arrange
		const invalidPath = "testdata/does_not_exist.json"

		// Act
		_, err := LoadConfig(invalidPath)

		// Assert
		if err == nil {
			t.Error("Expected error for invalid path but got none")
		}
	})

	t.Run("invalid JSON", func(t *testing.T) {
		// Arrange
		const invalidJSONPath = "testdata/invalid.json"

		// Act
		_, err := LoadConfig(invalidJSONPath)

		// Assert
		if err == nil {
			t.Error("Expected error for invalid JSON but got none")
		}
	})

}
