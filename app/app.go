package app

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/lpww/govs/pkg"
	"github.com/thatisuday/commando"
)

func Releases() error {
	fmt.Println("A full go release history can be found here: https://go.dev/doc/devel/release")
	return nil
}

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

	versions, err := getVersions()
	if err != nil {
		return err
	}

	if !pkg.BinExists("go") && len(versions) > 0 {
		fmt.Println("Warning: no default go version has been set. Use `govs list` to see installed vesions. Then set a default with `govs set <version>`. Ignore this message if you ran `govs get <version>`.")
	}

	return nil
}

func Install(args map[string]commando.ArgValue) error {
	v := args["version"].Value

	if pkg.BinExists("go") {
		return installVersion("go", v)
	}

	versions, err := getVersions()
	if err != nil {
		return err
	}

	if !pkg.BinExists("go") && len(versions) > 0 {
		defaultBin := fmt.Sprintf("go%s", versions[0]) // use the first installed go version as a default
		fmt.Printf("Warning: no default go version has been set. Using the existing %s to install go%s.\n", defaultBin, v)
		return installVersion(defaultBin, v)
	}

	if !pkg.BinExists("go") {

		// install temp go version if no bin found on the system
		fmt.Println("Warning: no go version found. Installing a temporary one.")

		// create a temp dir
		fmt.Println("Info: creating temp dir")
		tmp, err := os.MkdirTemp("", "govs")
		if err != nil {
			return errors.New(fmt.Sprintf("Error: temp dir could not be created.\n%s", err.Error()))
		}
		defer os.RemoveAll(tmp)

		// download go to tmp dir
		fmt.Println("Info: downloading temp go binary")
		tarGz := tmp + fmt.Sprintf("/go%s.tar.gz", v)
		goDownloadUrl := fmt.Sprintf("https://go.dev/dl/go%s.%s-%s.tar.gz", v, runtime.GOOS, runtime.GOARCH)
		if err = pkg.DownloadFile(tarGz, goDownloadUrl); err != nil {
			return errors.New(fmt.Sprintf("Error: go%s download failed.\n%s", v, err.Error()))
		}

		// ungzip downloaded go binary
		fmt.Println("Info: unpacking temp go binary")
		tar := fmt.Sprintf("%s/go%s.tar", tmp, v)
		if err := pkg.UnGzip(tarGz, tar); err != nil {
			return errors.New(fmt.Sprintf("Error: file, %s, could not be ungziped.\n%s", tarGz, err.Error()))
		}

		// untar downloaded go binary
		if err := pkg.Untar(tar, tmp); err != nil {
			return errors.New(fmt.Sprintf("Error: file, %s, could not be untared.\n%s", tar, err.Error()))
		}

		// use the temp go version to install a permanent go version
		fmt.Println("Info: installing requested go version")
		tmpBin := fmt.Sprintf("%s/go/bin/go", tmp)
		if err := installVersion(tmpBin, v); err != nil {
			return err
		}
	}

	return nil
}

func getVersions() (versions []string, err error) {
	src := pkg.GetSrcDir()
	entries, err := os.ReadDir(src)
	if err != nil {
		return versions, errors.New(fmt.Sprintf("Error: src directory, %s, could not be read.\n%s", src, err.Error()))
	}

	for _, e := range entries {
		version := pkg.TrimLeftChars(e.Name(), 2)
		versions = append(versions, version)
	}

	return versions, err
}

func List() error {
	versions, err := getVersions()
	if err != nil {
		return err
	}

	for _, v := range versions {
		fmt.Println(v)
	}

	return nil
}

func Set(args map[string]commando.ArgValue) error {
	v := args["version"].Value

	b := pkg.GetBinDir()
	s := pkg.GetSrcDir()

	vDir := fmt.Sprintf("%s/go%s", s, v)

	vGoBin := fmt.Sprintf("%s/bin/go", vDir)
	goBin := fmt.Sprintf("%s/go", b)

	if err := pkg.Symlink(vGoBin, goBin); err != nil {
		return fmt.Errorf("[pkg.Symlink]: %w", err)
	}

	vGofmtBin := fmt.Sprintf("%s/bin/gofmt", vDir)
	gofmtBin := fmt.Sprintf("%s/gofmt", b)

	if err := pkg.Symlink(vGofmtBin, gofmtBin); err != nil {
		return fmt.Errorf("[pkg.Symlink]: %w", err)
	}

	fmt.Printf("Success: the go command will now use %s.\n", v)

	return nil
}

func Remove(args map[string]commando.ArgValue) error {
	binDir := pkg.GetBinDir()
	srcDir := pkg.GetSrcDir()
	v := args["version"].Value

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

	fmt.Printf("Success: go%s has been removed.\n", v)

	versions, err := getVersions()
	if err != nil {
		return err
	}

	if !pkg.BinExists("go") && len(versions) == 0 {
		fmt.Println("Warning: you have removed the only go version from your system. Use `govs get <version>` to install and set a new default.")
	}

	if !pkg.BinExists("go") && len(versions) > 0 {
		fmt.Println("Warning: you have removed the default go version from your system. Use `govs list` to see installed vesions. Then set a new default with `govs set <version>`.")
	}

	return nil
}
