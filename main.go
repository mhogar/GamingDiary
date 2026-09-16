package main

import (
	"app/cmds"
	build_cmd "app/cmds/build"
	"fmt"
	"os"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/got-style/style"
)

func main() {
	runner := command.NewRunner(
		build_cmd.NewBuildCommand(),
		cmds.NewYoutubeCommand(),
		cmds.NewScriptsCommand(),
	)

	if len(os.Args) < 2 {
		runner.ListCommands()
		return
	}

	if err := runner.RunCommand(os.Args[1], os.Args[2:]); err != nil {
		fmt.Printf("%s %s\n", style.BoldError.Sprint("ERROR:"), err)
	}
}
