package main

import (
	"fmt"
	"govs/app"
	"io"
	"net/http"
	"os"

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

func FileExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func DirExists(path string) bool {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return info.IsDir()
}

func FatalError(message string) {
	fmt.Println(message)
	os.Exit(1)
}

func DownloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}

func main() {
	commando.
		SetExecutableName("govs").
		SetVersion("v1.0.0").
		SetDescription("A tool for installing and managing multiple go versions")

	commando.
		Register("install").
		SetDescription("This command installs the specified version of go and makes it available with a version specific binary, go<version>. Eg: `go1.18.5`").
		SetShortDescription("Install a go version").
		AddArgument("version", "The version to install", "").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			err := app.Install(args)
			if err != nil {
				FatalError(err.Error()) // todo: refactor to take an error and extract the string within fatalerror
			}
		})

	commando.
		Register("set").
		SetDescription("This command uses a symlink to create a default version for the `go` command").
		SetShortDescription("Set the default go version").
		AddArgument("version", "The version to set as the default", "").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			err := app.Set(args)
			if err != nil {
				FatalError(err.Error())
			}
		})

	commando.
		Register("remove").
		SetDescription("This command removes the specified version of go").
		SetShortDescription("Remove a go version").
		AddArgument("version", "The version to set as the default", "").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			d := GetDirs()
			v := args["version"].Value

			// todo: warn against removing currently set version

			vBin := fmt.Sprintf("%s/go%s", d.Bin, v)
			vSrc := fmt.Sprintf("%s/go%s", d.Src, v)

			if !FileExists(vBin) && !DirExists(vSrc) {
				FatalError(fmt.Sprintf("Error: go version %s is not installed. Please run `govs list` to see the installed versions", v))
			}

			if FileExists(vBin) {
				if err := os.Remove(vBin); err != nil {
					FatalError(fmt.Sprintf("Error: go version %s binary, %s, could not be removed.\n%s", v, vBin, err.Error()))
				}
			}

			if DirExists(vSrc) {
				if err := os.RemoveAll(vSrc); err != nil {
					FatalError(fmt.Sprintf("Error: go version %s src, %s, could not be removed.\n%s", v, vSrc, err.Error()))
				}
			}
		})

	commando.Parse(nil) // nil tells commando to use stdin arguments
}
