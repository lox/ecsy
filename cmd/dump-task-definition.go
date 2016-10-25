package cmd

import (
	"fmt"

	"github.com/lox/ecsy/api"
	"github.com/lox/ecsy/compose"
	"gopkg.in/alecthomas/kingpin.v2"
)

func ConfigureDumpTaskDefinition(app *kingpin.Application, svc api.Services) {
	var composeFiles []string

	cmd := app.Command("dump-task-definition", "Dump the task definition for a given set of docker-compose files")
	cmd.Alias("dump")

	cmd.Arg("files", "The docker-compose files to use").
		Required().
		ExistingFilesVar(&composeFiles)

	cmd.Action(func(c *kingpin.ParseContext) error {
		t := compose.Transformer{
			ComposeFiles: composeFiles,
		}

		taskDefinitionInput, err := t.Transform()
		if err != nil {
			return err
		}

		fmt.Println(taskDefinitionInput)

		return nil
	})
}
