package utils

import (
	"os/exec"
	"runtime"
)

func OpenURL(url string) error {
	switch runtime.GOOS {
	case "windows":
		return exec.Command("cmd", "/c", "start", url).Start()
	case "darwin":
		return exec.Command("open", url).Start()
	default: // Linux 或其他系统
		return exec.Command("xdg-open", url).Start()
	}
}
