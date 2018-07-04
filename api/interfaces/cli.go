package interfaces

import (
	"flag"
	"fmt"
	"io"
	"os"

	prompt "github.com/c-bata/go-prompt"
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
}

// Run invokes the CLI with the given arguments
func (c *CLI) Run(args []string) int {
	conf, err := c.loadConfig(args[1:])
	if err != nil {
		fmt.Fprintf(c.ErrStream, "args parse error: %v\n", err)
		return ExitCodeParseError
	}

	p := c.preparePrompt(conf)
	p.Run()

	// TODO: This is dead code. Invoke os.Exit by the prompt.Run
	return ExitCodeOK
}

func (c *CLI) loadConfig(args []string) (*config.Config, error) {
	conf, err := config.Load()
	if err != nil {
		return nil, err
	}

	// TODO: Implements usage
	// flag.Usage = func() {}
	flag.Parse()

	return conf, nil
}

func (c *CLI) preparePrompt(conf *config.Config) *prompt.Prompt {
	repository, err := bigtable.NewBigtableRepository(conf.Project, conf.Instance)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialized bigtable repository:%v", err)
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
	)
}
