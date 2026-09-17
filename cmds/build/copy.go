package build_cmd

import (
	"app/data"
	"app/util"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/binarysoupdev/go-extensions/errors"
	"github.com/binarysoupdev/go-extensions/file"
	"github.com/binarysoupdev/got-style/style"
)

type fileStats struct {
	Created  int
	UpToDate int
	NotFound int
}

func (s fileStats) Print() {
	s.printStat(style.Create, "created, ", s.Created)
	s.printStat(style.Info, "up-to-date, ", s.UpToDate)
	s.printStat(style.Error, "not found\n", s.NotFound)
}

func (s fileStats) printStat(style style.Style, text string, count int) {
	style.Printf("[%d] %s %s", count, util.Pluralize("file", "files", count), text)
}

//===========================================

func (cmd BuildCommand) copyFiles(dest, src string, stats *fileStats, files []string) error {
	for _, f := range files {
		if data.URL_REGEX.MatchString(f) {
			continue
		}

		if err := cmd.copyFileIfNewer(filepath.Join(dest, f), filepath.Join(src, f), stats); err != nil {
			return err
		}
	}
	return nil
}

func (cmd BuildCommand) copyFileIfNewer(dest, src string, stats *fileStats) error {
	newer, err := cmd.isFileNewer(src, dest)
	if err != nil {
		cmd.logger.Printf("[NOT FOUND] %s\n", src)
		stats.NotFound++
		return nil
	}

	if !newer {
		cmd.logger.Printf("[UP_TO_DATE] %s\n", src)
		stats.UpToDate++
		return nil
	}

	style.Create.Printf("+ %s -> %s\n", src, dest)

	if err := file.Copy(dest, src); err != nil {
		cmd.logError(err, fmt.Sprintf("%s -> %s", src, dest))
		panic(err)
	} else {
		cmd.logger.Printf("[COPIED] %s -> %s", src, dest)
		stats.Created++
	}
	return nil
}

func (cmd BuildCommand) isFileNewer(src, compare string) (bool, error) {
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
