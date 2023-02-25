package app

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/thatisuday/commando"
)

func BinExists(path string) bool {
	if _, err := exec.LookPath(path); err != nil {
		return false
	}
	return true
}

func Install(args map[string]commando.ArgValue) error {
	v := args["version"].Value

	if BinExists("go") {
		// run go install and go download

		vUrl := fmt.Sprintf("golang.org/dl/go%s@latest", v)
		install := exec.Command("go", "install", vUrl)
		install.Stdout = os.Stdout
		install.Stderr = os.Stderr
		if err := install.Run(); err != nil {
			return errors.New(fmt.Sprintf("Error: go%s could not be installed. Please ensure it is a valid version", v))
		}

		vBin := fmt.Sprintf("go%s", v)
		download := exec.Command(vBin, "download")
		download.Stdout = os.Stdout
		download.Stderr = os.Stderr
		if err := download.Run(); err == nil {
			return errors.New(fmt.Sprintf("Error: go%s could not be downloaded", v))
		}
	} else {
		// install go to temp dir
		// if tmp, err := os.MkdirTemp("", "govs"); err != nil {
		// 	FatalError(fmt.Sprintf("\nError: temp dir could not be created.\n%s", err.Error()))
		// }
		// defer os.RemoveAll(tmp)

		// // install go to temp dir
		// goDownloadUrl := "https://go.dev/dl/go1.20.1." + runtime.GOOS + "-" + runtime.GOARCH + ".tar.gz"
		// fmt.Println(goDownloadUrl)
		// err := DownloadFile(tmp+"/go1.20.1.tar.gz", goDownloadUrl)
		// if err != nil {
		// 	FatalError(fmt.Sprintf("\nError: go1.20.1 download failed.\n%s", err.Error()))
		// }
	}

	// unzip downloaded go binary

	// run the bin exists logic above

	// uninstall temp go

	// warn if no go binary exists in the path and recommend that the user run govs set v

	return nil
}
