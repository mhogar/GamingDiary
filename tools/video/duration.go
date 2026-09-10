package video

import (
	"os/exec"
	"strconv"

	"github.com/binarysoupdev/go-extensions/errors"
)

func CalcVideoDuration(path string) (float32, error) {
	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries", "format=duration", "-of", "default=noprint_wrappers=1:nokey=1", path)

	out, err := runCommand(cmd)
	if err != nil {
		return 0, err
	}

	f, err := strconv.ParseFloat(out, 32)
	if err != nil {
		return 0, errors.Chain(err, "error parsing duration")
	}

	return float32(f), nil
}
