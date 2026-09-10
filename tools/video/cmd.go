package video

import (
	"bytes"
	"os/exec"

	"github.com/binarysoupdev/go-extensions/errors"
)

func runCommand(cmd *exec.Cmd) (string, error) {
	stdout := bytes.Buffer{}
	cmd.Stdout = &stdout

	stderr := bytes.Buffer{}
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", errors.Chain(err, stderr.String())
	}
	return stdout.String()[:stdout.Len()-1], nil
}
