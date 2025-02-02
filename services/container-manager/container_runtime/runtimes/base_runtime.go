package runtime

import (
	"os"
)

type RuntimeBase struct {
	Name       string
	Runtime    string
	SocketPath string
}

func (runtime RuntimeBase) SocketExists() bool {
	_, err := os.Stat(runtime.SocketPath)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}
