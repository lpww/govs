package main

import (
	"fmt"
	"govs/app"
	"os"

	"github.com/thatisuday/commando"
)

func FatalError(message string) {
	fmt.Println(message)
	os.Exit(1)
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
			err := app.Remove(args)
			if err != nil {
				FatalError(err.Error())
			}
		})

	commando.Parse(nil) // nil tells commando to use stdin arguments
}
