package cmds

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/file"
	"github.com/binarysoupdev/go-extensions/json"
	"github.com/binarysoupdev/got-style/style"
)

type DeployCommand struct {
	command.CommandBase
	command.FlagCommand
}

func NewDeployCommand() *DeployCommand {
	return &DeployCommand{
		CommandBase: command.NewCommandBase("deploy", "copy static files"),
	}
}

func (cmd *DeployCommand) Initialize() error {
	cmd.InitFlagSet(cmd.Name, cmd.Description)
	return nil
}

func (cmd DeployCommand) Run(args []string) error {
	dest := cmd.Flags.String("dest", "", "the dest path")
	cmd.ParseFlags(args)

	if *dest == "" {
		return errors.New("\"dest\" cannot be empty")
	}

	data, err := json.UnmarshalFile[BaseData]("data/index.json")
	if err != nil {
		return errors.Chain(err, "error reading data file")
	}

	err = cmd.copyFiles(*dest, "public", []string{"index.html", "style.css", "script.js", "background.png"})
	if err != nil {
		return errors.Chain(err, "error copying root files")
	}

	for _, series := range data.Series {
		err := cmd.copySeries(*dest, series)
		if err != nil {
			return errors.Chain(err, fmt.Sprintf("error copying \"%s\" files", series))
		}
	}

	return nil
}

func (cmd DeployCommand) copySeries(dest, series string) error {
	data, err := json.UnmarshalFile[SeriesData](filepath.Join("data", series, "index.json"))
	if err != nil {
		return errors.Chain(err, "error reading data file")
	}

	dest = filepath.Join(dest, series)

	files := []string{"index.html"}
	files = append(files, data.Stylesheets...)
	files = append(files, data.Background, data.Thumbnail)

	err = os.MkdirAll(dest, 0755)
	if err != nil {
		return errors.Chain(err, "error creating dest directory")
	}

	err = cmd.copyFiles(dest, filepath.Join("public", series), files)
	if err != nil {
		return errors.Chain(err, "error copying files")
	}

	return nil
}

func (cmd DeployCommand) copyFiles(dest, src string, files []string) error {
	errs := errors.Errors{}

	for _, f := range files {
		if err := cmd.copyFileIfNewer(filepath.Join(dest, f), filepath.Join(src, f)); err != nil {
			errs.Add(err)
		}
	}

	return errs.Collapse(", ")
}

func (cmd DeployCommand) copyFileIfNewer(dest, src string) error {
	style.Info.Print(src)

	newer, err := cmd.isFileNewer(src, dest)
	if err != nil {
		return err
	}

	if !newer {
		style.Info.Println(" -> up-to-date")
		return nil
	}

	err = file.Copy(dest, src)
	if err != nil {
		style.Error.Printf(" -> %s\n", dest)
		return errors.Chain(err, "error copying file")
	}

	style.Success.Printf(" -> %s\n", dest)
	return nil
}

func (cmd DeployCommand) isFileNewer(src, compare string) (bool, error) {
	stat, err := os.Stat(src)
	if err != nil {
		return false, errors.Chain(err, "error reading source file")
	}

	target, err := os.Stat(compare)
	if err != nil {
		return true, nil
	}

	return target.ModTime().Before(stat.ModTime()), nil
}
