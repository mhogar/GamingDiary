package build_cmd

import (
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
	Invalid  int
}

func (s fileStats) Print() {
	style.Create.Printf("[%d] files created, ", s.Created)
	style.Info.Printf("[%d] files up-to-date, ", s.UpToDate)
	style.Error.Printf("[%d] files not found, ", s.NotFound)
	style.Error.Printf("[%d] files invalid\n", s.Invalid)
}

//===========================================

func (cmd BuildCommand) copyFiles(dest, src string, stats fileStats, files []string) {
	for _, f := range files {
		cmd.copyFileIfNewer(filepath.Join(dest, f), filepath.Join(src, f), &stats)
	}
	stats.Print()
}

func (cmd BuildCommand) copyFileIfNewer(dest, src string, stats *fileStats) {
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

	style.Create.Printf("+ %s -> %s\n", src, dest)

	if err := file.Copy(dest, src); err != nil {
		cmd.logError(err, fmt.Sprintf("%s -> %s", src, dest))
		stats.Invalid++
	} else {
		cmd.logger.Printf("[COPIED] %s -> %s", src, dest)
		stats.Created++
	}
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
