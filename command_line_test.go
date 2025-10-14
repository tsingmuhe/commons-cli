package cli_test

import (
	"os"
	"testing"

	cli "github.com/tsingmuhe/commons-cli"
)

type RootCommand struct {
	Version bool      `short:"v" long:"--version" description:"Show version"`
	Help    bool      `short:"h" long:"--help" description:"Show help"`
	File    [5]string `description:"Specify a file"`

	*SubCommand
}

func (r *RootCommand) Name() string {
	return "sun"
}

func (r *RootCommand) Description() string {
	return "sun root command"
}

func (r *RootCommand) Run() int {
	return 0
}

type SubCommand struct {
}

func (s *SubCommand) Name() string {
	return "chp"
}

func (s *SubCommand) Description() string {
	return "sun sub command"
}

func (s *SubCommand) Run() int {
	return 0
}

func TestCommandLine_New(t *testing.T) {
	cmdLine, err := cli.Create(new(RootCommand), "1.0.0")
	if err != nil {
		return
	}

	_ = cmdLine.Run(os.Args)
}
