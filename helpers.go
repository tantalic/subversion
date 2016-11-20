package subversion

import (
	"bytes"
	"os/exec"
)

var (
	newlineCharacter = []byte("\n")
	emptyCharacter   = []byte{}
)

func cleanCmdOutput(cmd *exec.Cmd) (string, error) {
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}

	out = bytes.Replace(out, newlineCharacter, emptyCharacter, 1)
	return string(out), nil
}
