package interfaces

import (
	"flag"
	"fmt"
	"io"
	"os"

	prompt "github.com/c-bata/go-prompt"
	"github.com/pkg/errors"
	"github.com/takashabe/btcli/api/application"
	"github.com/takashabe/btcli/api/config"
	"github.com/takashabe/btcli/api/infrastructure/bigtable"
)

// exit codes
const (
	ExitCodeOK = 0

	// Specific error codes. begin 10-
	ExitCodeError = 10 + iota
	ExitCodeParseError
	ExitCodeInvalidArgsError
)

// CLI is the command line interface object
type CLI struct {
	OutStream io.Writer
	ErrStream io.Writer

	Version string
}

// Run invokes the CLI with the given arguments
func (c *CLI) Run(args []string) int {
	conf, err := c.loadConfig(args)
	if err != nil {
		fmt.Fprintf(c.ErrStream, "args parse error: %v\n", err)
		return ExitCodeParseError
	}

	p, err := c.preparePrompt(conf)
	if err != nil {
		fmt.Fprintf(c.ErrStream, "failed to initialized prompt: %v\n", err)
		return ExitCodeError
	}

	fmt.Fprintf(c.OutStream, "Version: %s\n", c.Version)
	fmt.Fprintf(c.OutStream, "Please use `exit` or `Ctrl-D` to exit this program.\n")
	p.Run()

	// TODO: This is dead code. Invoke os.Exit by the prompt.Run
	return ExitCodeOK
}

func (c *CLI) loadConfig(args []string) (*config.Config, error) {
	conf := config.NewConfig(c.ErrStream)
	err := conf.Load()
	if err != nil {
		return nil, err
	}

	flag.Usage = func() {
		usage(c.OutStream)
	}
	return conf, nil
}

func usage(w io.Writer) {
	fmt.Fprintf(w, "Usage: %s [flags] <command> ...\n", os.Args[0])
	flag.CommandLine.SetOutput(w)
	flag.CommandLine.PrintDefaults()
}

func (c *CLI) preparePrompt(conf *config.Config) (*prompt.Prompt, error) {
	repository, err := bigtable.NewBigtableRepository(conf.Project, conf.Instance)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to initialized bigtable repository:%v", err)
	}
	tableInteractor := application.NewTableInteractor(repository)
	rowsInteractor := application.NewRowsInteractor(repository)

	executor := Executor{
		outStream:       c.OutStream,
		errStream:       c.ErrStream,
		rowsInteractor:  rowsInteractor,
		tableInteractor: tableInteractor,
	}
	completer := Completer{
		tableInteractor: tableInteractor,
	}

	return prompt.New(
		executor.Do,
		completer.Do,
		// TODO: Add histories from the history file.
		// prompt.OptionHistory(),
	), nil
}
