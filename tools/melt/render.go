package melt

import (
	"fmt"
	"os/exec"
)

func RenderVideo(path string) {
	cmd := exec.Command("notify-send", fmt.Sprintf("Video Path: %s", path))
	cmd.Run()
}
