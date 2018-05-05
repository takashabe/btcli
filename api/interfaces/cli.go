package interfaces

import (
	"fmt"
	"os"

	"github.com/takashabe/btcli/api/application"
	"github.com/takashabe/btcli/api/infrastructure/bigtable"
)

var (
	tableInteractor *application.TableInteractor
	rowsInteractor  *application.RowsInteractor
)

// TODO: delegate to the main package
func init() {
	var (
		project  = "test-project"
		instance = "test-instance"
	)

	repository, err := bigtable.NewBigtableRepository(project, instance)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialized bigtable repository:%v", err)
	}
	tableInteractor = application.NewTableInteractor(repository)
	rowsInteractor = application.NewRowsInteractor(repository)
}
