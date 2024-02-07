package internal

import (
	"encoding/hex"
	"fmt"
	"os"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func LogError(errLog error) {
	errMsg := errLog.Error()
	logs, err := os.OpenFile("logs", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to log error: " + err.Error())
	}
	defer logs.Close()
	_, err = logs.WriteString(time.Now().Format("2006-01-02 15:04:05") + "|" + errMsg + "\n")
	if err != nil {
		fmt.Println("Failed to log error: " + err.Error())
	}
}

func LogString(msg string) {
	logs, err := os.OpenFile("logs", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Println("Failed to log message: " + err.Error())
	}
	defer logs.Close()
	_, err = logs.WriteString(time.Now().Format("2006-01-02 15:04:05") + "|" + msg + "\n")
	if err != nil {
		fmt.Println("Failed to log message: " + err.Error())
	}
}

func hexToRaylib(color string) rl.Color {
	rgba, err := hex.DecodeString(color[1:])
	if err != nil {
		fmt.Println("Unable to parse color")
		return rl.White
	}
	return rl.NewColor(rgba[0], rgba[1], rgba[2], 255)
}

func absInt(number int) int {
	if number < 0 {
		return number * -1
	}
	return number
}

func IsUpdateTick(s State) bool {
	currentTime := time.Now().UnixMilli()
	if currentTime-int64(s.lastUpdateTime) > 200 {
		return true
	}
	return false
}
