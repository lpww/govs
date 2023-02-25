package app

import (
	"errors"
	"fmt"
	"govs/pkg"
	"os"
	"os/exec"
	"runtime"

	"github.com/thatisuday/commando"
)

func installVersion(cmd string, version string) error {
	vUrl := fmt.Sprintf("golang.org/dl/go%s@latest", version)
	install := exec.Command(cmd, "install", vUrl)
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
		return installVersion("go", v)
	} else {
		// todo: if they have go versions installed but not set, warn them and ask to set a go vesion before continuing

		// todo: refactor to manually install the specified version instead of using a temp go version
		// 1: extract the tarball to $HOME/sdk (pkg.GetSrcDir())
		// 2. rename $HOME/sdk/go to $HOME/sdk/go<version>
		// 3. copy $HOME/sdk/go<version> to pkg.GetBinDir()/go<version>

		// install temp go version if no bin found on the system
		fmt.Println("Warning: no go binary detected. Installing a temporary one.")

		// create a temp dir
		fmt.Println("Info: creating temp dir")
		tmp, err := os.MkdirTemp("", "govs")
		if err != nil {
			return errors.New(fmt.Sprintf("Error: temp dir could not be created.\n%s", err.Error()))
		}
		defer os.RemoveAll(tmp)

		// download go to tmp dir
		fmt.Println("Info: downloading temp go binary")
		tarGz := tmp + "/go1.20.1.tar.gz"
		goDownloadUrl := "https://go.dev/dl/go1.20.1." + runtime.GOOS + "-" + runtime.GOARCH + ".tar.gz"
		if err = pkg.DownloadFile(tarGz, goDownloadUrl); err != nil {
			return errors.New(fmt.Sprintf("Error: go1.20.1 download failed.\n%s", err.Error()))
		}

		// ungzip downloaded go binary
		fmt.Println("Info: ungziping temp go binary")
		tar := tmp + "/go1.20.1.tar"
		if err := pkg.UnGzip(tarGz, tar); err != nil {
			return errors.New(fmt.Sprintf("Error: file, %s, could not be ungziped.\n%s", tarGz, err.Error()))
		}

		// untar downloaded go binary
		fmt.Println("Info: untaring temp go binary")
		src := tmp
		if err := pkg.Untar(tar, src); err != nil {
			return errors.New(fmt.Sprintf("Error: file, %s, could not be untared.\n%s", tar, err.Error()))
		}

		// use the temp go version to install a permanent go version
		fmt.Println("Info: installing specified go version")
		tmpBin := src + "/go/bin/go"
		if err := installVersion(tmpBin, v); err != nil {
			return err
		}

		// warn if no go binary exists in the path and recommend that the user run govs set v
	}

	return nil
}

func List() error {
	src := pkg.GetSrcDir()
	versions, err := os.ReadDir(src)
	if err != nil {
		return errors.New(fmt.Sprintf("Error: src directory, %s, could not be read.\n%s", src, err.Error()))
	}

	fmt.Println("Installed go versions:")
	for _, v := range versions {
		fmt.Println(pkg.TrimLeftChars(v.Name(), 2))
	}

	return nil
}

func Set(args map[string]commando.ArgValue) error {
	d := pkg.GetBinDir()
	v := args["version"].Value

	vBin := fmt.Sprintf("%s/go%s", d, v)
	goBin := fmt.Sprintf("%s/go", d)

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
	binDir := pkg.GetBinDir()
	srcDir := pkg.GetSrcDir()
	v := args["version"].Value

	// todo: warn against removing currently set version

	vBin := fmt.Sprintf("%s/go%s", binDir, v)
	vSrc := fmt.Sprintf("%s/go%s", srcDir, v)

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
