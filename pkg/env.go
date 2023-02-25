package pkg

import "os"

func GetBinDir() string {
	v, ok := os.LookupEnv("GOPATH")

	if !ok || v == "" {
		// use a default if $GOPATH is not set
		return os.ExpandEnv("$HOME/go/bin")
	}

	return os.ExpandEnv("$GOPATH/bin")
}

func GetSrcDir() string {
	return os.ExpandEnv("$HOME/sdk")
}
