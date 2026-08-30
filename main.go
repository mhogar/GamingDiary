package main

/*/media/ryan/Seagate Portable Drive/Videos/Gaming Diary/public*/

import (
	"fmt"
	"local/cmds"
	"os"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/got-style/style"
)

func main() {
	runner := command.NewRunner(
		cmds.NewRenderCommand(),
		cmds.NewRenameCommand(),
		cmds.NewVideoStatCommand(),
		cmds.NewDeployCommand(),
	)

	if len(os.Args) < 2 {
		runner.ListCommands()
		return
	}

	if err := runner.RunCommand(os.Args[1], os.Args[2:]); err != nil {
		fmt.Printf("%s %s\n", style.BoldError.Sprint("ERROR:"), err)
	}
}
