package cmds

import (
	"fmt"
	"local/data"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/binarysoupdev/go-commando/command"
	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/file"
	"github.com/binarysoupdev/go-extensions/json"
	"github.com/binarysoupdev/got-style/style"
)

type DeployStats struct {
	Copied   int
	UpToDate int
	NotFound int
}

type DeployCommand struct {
	command.CommandBase
	command.FlagCommand

	logger *log.Logger
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
	public := cmd.Flags.String("public", "", "the public path")
	cmd.ParseFlags(args)

	if *public == "" {
		return errors.New("\"public\" cannot be empty")
	}

	err := os.MkdirAll(*public, 0755)
	if err != nil {
		return errors.Chain(err, "error creating public directory")
	}

	f, err := os.Create(fmt.Sprintf("logs/deploy-%s.txt", time.Now().Format(time.DateTime)))
	if err != nil {
		return errors.Chain(err, "error creating log file")
	}
	defer f.Close()
	cmd.logger = log.New(f, "", log.Ltime)

	root, err := json.UnmarshalFile[data.Root]("data/index.json")
	if err != nil {
		return errors.Chain(err, "error reading data file")
	}

	style.Bold.Println("root")
	cmd.copyFiles(*public, "public", []string{"index.html", "style.css", "script.js", "background.png"})

	for _, series := range root.Series {
		err := cmd.copySeries(*public, series)
		if err != nil {
			return errors.Chain(err, fmt.Sprintf("error copying \"%s\" files", series))
		}
	}

	return nil
}

func (cmd DeployCommand) copySeries(dest, name string) error {
	dest = filepath.Join(dest, name)

	series, err := json.UnmarshalFile[data.Series](filepath.Join("data", name, "index.json"))
	if err != nil {
		return errors.Chain(err, "error reading data file")
	}

	entires, err := json.UnmarshalFile[data.Entries](filepath.Join("data", name, series.Entries))
	if err != nil {
		return errors.Chain(err, "error reading entires file")
	}

	err = os.MkdirAll(dest, 0755)
	if err != nil {
		return errors.Chain(err, "error creating dest directory")
	}

	files := make([]string, 0, 3+len(series.Stylesheets)+(len(entires.Entries)*2))
	files = append(files, "index.html", series.Background, series.Thumbnail)
	files = append(files, series.Stylesheets...)

	for _, entry := range entires.Entries {
		files = append(files, entry.Thumbnail, entry.Video)
	}

	style.Bold.Println(name)
	cmd.copyFiles(dest, filepath.Join("public", name), files)
	return nil
}

func (cmd DeployCommand) copyFiles(dest, src string, files []string) {
	stats := DeployStats{}

	for _, f := range files {
		cmd.copyFileIfNewer(filepath.Join(dest, f), filepath.Join(src, f), &stats)
	}
	fmt.Println("\n---")

	style.Create.Printf("[%d] files copied, ", stats.Copied)
	style.Info.Printf("[%d] file up-to-date, ", stats.UpToDate)
	style.Error.Printf("[%d] files not found\n", stats.NotFound)
}

func (cmd DeployCommand) copyFileIfNewer(dest, src string, stats *DeployStats) {
	newer, err := cmd.isFileNewer(src, dest)
	if err != nil {
		cmd.logger.Printf("[NOT FOUND] %s\n", src)
		stats.NotFound++
		return
	}

	if !newer {
		cmd.logger.Printf("[UP_TO_DATE] %s\n", src)
		stats.UpToDate++
		return
	}

	style.Success.Printf("\r+ %s ", dest)

	if err := file.Copy(dest, src); err != nil {
		cmd.logger.Printf("[ERROR] %s -> %s\n", src, dest)
		cmd.logger.Println(err)
		return
	}

	cmd.logger.Printf("[COPIED] %s | %s", src, dest)
	stats.Copied++
}

func (cmd DeployCommand) isFileNewer(src, compare string) (bool, error) {
	stat, err := os.Stat(src)
	if err != nil {
		return false, errors.Chain(err, "error reading source file")
	}
	srcTime := stat.ModTime().Truncate(time.Second)

	stat, err = os.Stat(compare)
	if err != nil {
		return true, nil
	}
	destTime := stat.ModTime().Truncate(time.Second)

	return destTime.Before(srcTime), nil
}
