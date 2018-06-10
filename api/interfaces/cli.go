package interfaces

import (
	"flag"
	"fmt"
	"io"
	"os"

	prompt "github.com/c-bata/go-prompt"
	"github.com/pkg/errors"
	"github.com/takashabe/btcli/api/application"
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

const (
	defaultProject  = "test-project"
	defaultInstance = "test-instance"
)

// TODO: need require/optional params
type param struct {
	project  string
	instance string
}

// CLI is the command line interface object
type CLI struct {
	OutStream io.Writer
	ErrStream io.Writer
}

// Run invokes the CLI with the given arguments
func (c *CLI) Run(args []string) int {
	param := &param{}
	err := c.parseArgs(args[1:], param)
	if err != nil {
		fmt.Fprintf(c.ErrStream, "args parse error: %v\n", err)
		return ExitCodeParseError
	}

	p := c.preparePrompt(param)
	p.Run()

	// TODO: This is dead code. Invoke os.Exit by the prompt.Run
	return ExitCodeOK
}

func (c *CLI) parseArgs(args []string, p *param) error {
	flags := flag.NewFlagSet("param", flag.ContinueOnError)
	flags.SetOutput(c.ErrStream)

	flags.StringVar(&p.project, "project", defaultProject, `project ID, if unset uses gcloud configured project (default "test-project")`)
	flags.StringVar(&p.instance, "instance", defaultInstance, `Cloud Bigtable instance (default "test-instance")`)

	err := flags.Parse(args)
	if err != nil {
		return errors.Wrapf(err, "failed to parsed args")
	}
	return nil
}

func (c *CLI) preparePrompt(p *param) *prompt.Prompt {
	repository, err := bigtable.NewBigtableRepository(p.project, p.instance)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialized bigtable repository:%v", err)
	}
	tableInteractor := application.NewTableInteractor(repository)
	rowsInteractor := application.NewRowsInteractor(repository)

	executor := Executor{
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
