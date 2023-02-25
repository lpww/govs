package app

import (
	"errors"
	"fmt"
	"govs/pkg"
	"os"
	"os/exec"

	"github.com/thatisuday/commando"
)

type Env struct {
	Key   string
	Value string
}

var GOPATH = Env{"GOPATH", ""}

var HOME = Env{"HOME", ""}

func GetEnv(e Env, defaultValue string) Env {
	v, ok := os.LookupEnv(e.Key)

	if !ok || v == "" {
		fmt.Printf("Warning: $%s not set. Using default value of %s\n", e.Key, defaultValue)
		return Env{e.Key, defaultValue}
	}

	return Env{e.Key, v}
}

type Dirs struct {
	Bin string
	Src string
}

// todo: remove this fn and use ExpandEnv inline
func GetDirs() Dirs {
	home := GetEnv(HOME, "")
	goPath := GetEnv(GOPATH, fmt.Sprintf("%s/go", home.Value))

	return Dirs{
		Bin: fmt.Sprintf("%s/bin", goPath.Value),
		Src: fmt.Sprintf("%s/sdk", home.Value),
	}
}

func installVersion(version string) error {
	vUrl := fmt.Sprintf("golang.org/dl/go%s@latest", version)
	install := exec.Command("go", "install", vUrl)
	install.Stdout = os.Stdout
	install.Stderr = os.Stderr
	if err := install.Run(); err != nil {
		return errors.New(fmt.Sprintf("Error: go%s could not be installed. Please ensure it is a valid version\n%s", version, err.Error()))
	}

	vBin := fmt.Sprintf("go%s", version)
	download := exec.Command(vBin, "download")
	download.Stdout = os.Stdout
	download.Stderr = os.Stderr
	if err := download.Run(); err != nil {
		return errors.New(fmt.Sprintf("Error: go%s could not be downloaded\n%s", version, err.Error()))
	}

	return nil
}

func Install(args map[string]commando.ArgValue) error {
	v := args["version"].Value

	if pkg.BinExists("go") {
		return installVersion(v)

	} else {
		// install go to temp dir
		// if tmp, err := os.MkdirTemp("", "govs"); err != nil {
		// 	FatalError(fmt.Sprintf("\nError: temp dir could not be created.\n%s", err.Error()))
		// }
		// defer os.RemoveAll(tmp)

		// // install go to temp dir
		// goDownloadUrl := "https://go.dev/dl/go1.20.1." + runtime.GOOS + "-" + runtime.GOARCH + ".tar.gz"
		// fmt.Println(goDownloadUrl)
		// err := pkg.DownloadFile(tmp+"/go1.20.1.tar.gz", goDownloadUrl)
		// if err != nil {
		// 	FatalError(fmt.Sprintf("\nError: go1.20.1 download failed.\n%s", err.Error()))
		// }
	}

	// unzip downloaded go binary

	// run the bin exists logic above
	if err := installVersion(v); err != nil {
		return err
	}

	// uninstall temp go

	// warn if no go binary exists in the path and recommend that the user run govs set v

	return nil
}

func Set(args map[string]commando.ArgValue) error {
	d := GetDirs()
	v := args["version"].Value

	vBin := fmt.Sprintf("%s/go%s", d.Bin, v)
	goBin := fmt.Sprintf("%s/go", d.Bin)

	// todo: warn if the goroot is not $HOME/sdk/go* - why?

	if !pkg.FileExists(vBin) {
		return errors.New(fmt.Sprintf("Error: go version %s is not installed. Please run `govs install %s` and try again", v, v))
	}

	if _, err := os.Lstat(goBin); err == nil {
		if err := os.Remove(goBin); err != nil {
			return errors.New(fmt.Sprintf("Error: existing go binary, %s, could not be removed.\n%s", goBin, err.Error()))
		}
	}

	if err := os.Symlink(vBin, goBin); err != nil {
		return errors.New(fmt.Sprintf("Error: the default go version could not be set to %s.\n%s", v, err.Error()))
	}

	return nil
}

func Remove(args map[string]commando.ArgValue) error {
	d := GetDirs()
	v := args["version"].Value

	// todo: warn against removing currently set version

	vBin := fmt.Sprintf("%s/go%s", d.Bin, v)
	vSrc := fmt.Sprintf("%s/go%s", d.Src, v)

	if !pkg.FileExists(vBin) && !pkg.DirExists(vSrc) {
		return errors.New(fmt.Sprintf("Error: go version %s is not installed. Please run `govs list` to see the installed versions", v))
	}

	if pkg.FileExists(vBin) {
		if err := os.Remove(vBin); err != nil {
			return errors.New(fmt.Sprintf("Error: go version %s binary, %s, could not be removed.\n%s", v, vBin, err.Error()))
		}
	}

	if pkg.DirExists(vSrc) {
		if err := os.RemoveAll(vSrc); err != nil {
			return errors.New(fmt.Sprintf("Error: go version %s src, %s, could not be removed.\n%s", v, vSrc, err.Error()))
		}
	}

	return nil
}
