package build_cmd

func (cmd BuildCommand) logError(err error, msg string) {
	cmd.logger.Printf("[ERROR] %s\n  %s\n", msg, err)
}

func (cmd BuildCommand) logCreate(file string) {
	cmd.logger.Printf("[CREATE] %s\n", file)
}

func (cmd BuildCommand) logBuild(series string) {
	cmd.logger.Printf("[BUILD] %s\n", series)
}
