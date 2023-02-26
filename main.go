package main

import (
	"fmt"
	"os"

	"github.com/lpww/govs/app"

	"github.com/thatisuday/commando"
)

func FatalError(err error) {
	fmt.Println(err.Error())
	os.Exit(1)
}

func main() {
	commando.
		SetExecutableName("govs").
		SetVersion("v1.0.1").
		SetDescription("A tool for installing and managing multiple go versions")

	commando.
		Register("releases").
		SetDescription("This command displays all official go language releases available for download").
		SetShortDescription("List available go releases").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			err := app.Releases()
			if err != nil {
				FatalError(err)
			}
		})

	commando.
		Register("install").
		SetDescription("This command installs the specified version of go and makes it available with a version specific binary, go<version>. Eg: `go1.18.5`").
		SetShortDescription("Install a go version").
		AddArgument("version", "The version to install", "").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			err := app.Install(args)
			if err != nil {
				FatalError(err)
			}
		})

	commando.
		Register("list").
		SetDescription("This command lists the installed go versions available on the system").
		SetShortDescription("List installed go versions").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			err := app.List()
			if err != nil {
				FatalError(err)
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
				FatalError(err)
			}
		})

	commando.
		Register("get").
		SetDescription("This command will install the specified version and set it as the default version for the `go` command").
		SetShortDescription("Install and set the default go version").
		AddArgument("version", "The version to install and set as the default", "").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			if err := app.Install(args); err != nil {
				FatalError(err)
			}

			if err := app.Set(args); err != nil {
				FatalError(err)
			}
		})

	commando.
		Register("remove").
		SetDescription("This command removes the specified version of go").
		SetShortDescription("Remove a go version").
		AddArgument("version", "The version to remove", "").
		SetAction(func(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
			err := app.Remove(args)
			if err != nil {
				FatalError(err)
			}
		})

	commando.Parse(nil) // nil tells commando to use stdin arguments
}
