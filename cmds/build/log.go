package build_cmd

import (
	"strings"

	"github.com/binarysoupdev/got-style/style"
)

func (cmd BuildCommand) logError(err error, msg string) {
	cmd.logger.Printf("[ERROR] %s\n  %s\n", msg, err)
}

func (cmd BuildCommand) logCreate(file string) {
	cmd.logger.Printf("[CREATE] %s\n", file)
}

func (cmd BuildCommand) logBuild(series string) {
	cmd.logger.Printf("[BUILD] %s\n", strings.ToUpper(series))
}

func (cmd BuildCommand) printSeriesHeader(series string) {
	style.New(style.BOLD, style.UNDERLINE).Println(series)
}

func (cmd BuildCommand) printError(err string) {
	style.Error.Printf("x %s\n", err)
}

func (cmd BuildCommand) printCreate(file string) {
	style.Create.Printf("+ %s\n", file)
}
