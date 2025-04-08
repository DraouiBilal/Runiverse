package runtime

import (
	"os"
)

type RuntimeBase struct {
	Name       string
	Runtime    string
	SocketPath string
}

func (runtime RuntimeBase) SocketExists() (bool, error) {
	_, err := os.Stat(runtime.SocketPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
