package pkg

import (
	"errors"
	"fmt"
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

func Symlink(oldname string, newname string) error {
	if _, err := os.Lstat(newname); err == nil {
		if err := os.Remove(newname); err != nil {
			return errors.New(fmt.Sprintf("Error: existing binary, %s, could not be removed.\n%s", newname, err.Error()))
		}
	}

	if err := os.Symlink(oldname, newname); err != nil {
		return errors.New(fmt.Sprintf("Error: the default go version could not be set to %s.\n%s", oldname, err.Error()))
	}

	return nil
}
