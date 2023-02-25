package pkg

import (
	"os"
	"os/exec"
)

func BinExists(path string) bool {
	if _, err := exec.LookPath(path); err != nil {
		return false
	}
	return true
}

func DirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
