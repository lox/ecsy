package cmd

import (
	"fmt"
	"os"

	"github.com/lox/ecsy/api"
	"github.com/lox/ecsy/taskdefinitions"
	"gopkg.in/alecthomas/kingpin.v2"
)

func ConfigureDumpTaskDefinition(app *kingpin.Application, svc api.Services) {
	var file string

	cmd := app.Command("dump-task-definition", "Dump the task definition after environment interpolation")
	cmd.Alias("dump")

	cmd.Arg("file", "The task-definition file to use").
		Required().
		ExistingFileVar(&file)

	cmd.Action(func(c *kingpin.ParseContext) error {
		taskDefinitionInput, err := taskdefinitions.ParseFile(file, os.Environ())
		if err != nil {
			return err
		}

		fmt.Println(taskDefinitionInput.String())
		return nil
	})
}
